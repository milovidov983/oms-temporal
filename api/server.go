package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bojanz/httpx"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/milovidov983/oms-temporal/internal/cartorder"
	"github.com/milovidov983/oms-temporal/internal/signals"
	"github.com/milovidov983/oms-temporal/internal/signals/channels"
	"github.com/milovidov983/oms-temporal/internal/signals/routes"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"go.temporal.io/sdk/client"
)

var (
	HTTPPort = "8888" //os.Getenv("PORT")
	temporal client.Client
)

func main() {
	var err error
	temporal, err = client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	log.Println("Temporal client connected")

	r := mux.NewRouter()
	r.Handle("/order", http.HandlerFunc(CreateOrderHandler)).Methods("POST")
	r.Handle("/order/{workflowID}", http.HandlerFunc(GetOrderHandler)).Methods("GET")
	r.Handle("/order/{workflowID}/assembly", http.HandlerFunc(CompleteAssemblyHandler)).Methods("PUT")
	r.Handle("/order/{workflowID}/assembly-comment", http.HandlerFunc(ChangeAssemblyCommentHandler)).Methods("PUT")
	r.Handle("/cart/{workflowID}/delivery", http.HandlerFunc(CompleteDeliveryHandler)).Methods("PUT")
	r.Handle("/cart/{workflowID}/delivery-comment", http.HandlerFunc(ChangeDeliveryCommentHandler)).Methods("PUT")
	r.Handle("/cart/{workflowID}/cancel", http.HandlerFunc(CancelOrderHandler)).Methods("PUT")

	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	var cors = handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))

	http.Handle("/", cors(r))
	server := httpx.NewServer(":"+HTTPPort, http.DefaultServeMux)
	server.WriteTimeout = time.Second * 240

	log.Println("Starting server on port: " + HTTPPort)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := fmt.Sprintf("%d", time.Now().Unix())
	workflowID := "ORDER-" + orderID

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "ORDER_TASK_QUEUE",
	}

	order := models.OrderState{
		OrderID: orderID,
		Status:  models.OrderStatusCreated,
		// example
		Ordered:   []models.OrderLines{models.OrderLines{ProductID: 42, Quantity: 1, Price: 100}},
		Collected: make([]models.OrderLines, 0),
		Delivered: make([]models.OrderLines, 0),
	}

	we, err := temporal.ExecuteWorkflow(r.Context(), options, cartorder.CartOrderWorkflow, order)
	if err != nil {
		WriteError(w, err)
		return
	}

	res := make(map[string]interface{})
	res["order"] = order
	res["workflowID"] = we.GetID()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := temporal.QueryWorkflow(context.Background(), vars["workflowID"], "", "getOrder")
	if err != nil {
		WriteError(w, err)
		return
	}
	var res interface{}
	if err := response.Get(&res); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func CompleteAssemblyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var collected []models.OrderLines
	if err := json.NewDecoder(r.Body).Decode(&collected); err != nil {
		WriteError(w, err)
		return
	}

	// TODO:
	// При завершении этапа сборки необходимо сделать несколько синхронных
	// вызовов, чтобы оповестить сборщика о том что все ок. А именно
	// в собранных позициях все ок, оплата прошла корректно
	// 1. Проверить позиции собранного
	// 2. Произвести оплату

	// 3. Отправка сигнала в workflow заказа для дальнейшей обработки
	update := signals.SignalPayloadCompleteAssembly{
		Route:     routes.RouteTypeCompleteAssembly,
		Collected: collected,
	}
	signalName := channels.SignalNameCompleteAssemblyChannel
	workflowID := vars["workflowID"]

	err := temporal.SignalWorkflow(r.Context(), workflowID, "", signalName, update)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func ChangeAssemblyCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var requestBody models.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		WriteError(w, err)
		return
	}

	update := signals.SignalPayloadChangeAssemblyComment{
		Route:   routes.RouteTypeChangeAssemblyComment,
		Comment: requestBody.Comment,
	}
	signalName := channels.SignalNameChangeAssemblyCommentChannel
	workflowID := vars["workflowID"]

	err := temporal.SignalWorkflow(r.Context(), workflowID, "", signalName, update)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func CompleteDeliveryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var delivered []models.OrderLines
	if err := json.NewDecoder(r.Body).Decode(&delivered); err != nil {
		WriteError(w, err)
		return
	}

	update := signals.SignalPayloadCompleteDelivery{
		Route:     routes.RouteTypeCompleteDelivery,
		Delivered: delivered,
	}
	signalName := channels.SignalNameCompleteDeliveryChannel
	workflowID := vars["workflowID"]

	err := temporal.SignalWorkflow(r.Context(), workflowID, "", signalName, update)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func ChangeDeliveryCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var requestBody models.CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		WriteError(w, err)
		return
	}

	update := signals.SignalPayloadChangeDeliveryComment{
		Route:   routes.RouteTypeChangeDeliveryComment,
		Comment: requestBody.Comment,
	}
	signalName := channels.SignalNameChangeDeliveryCommentChannel
	workflowID := vars["workflowID"]

	err := temporal.SignalWorkflow(r.Context(), workflowID, "", signalName, update)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var requestBody models.ReasonRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		WriteError(w, err)
		return
	}

	update := signals.SignalPayloadCancelOrder{
		Route:  routes.RouteTypeCancelOrder,
		Reason: requestBody.Reason,
	}
	signalName := channels.SignalNameCancelOrderChannel
	workflowID := vars["workflowID"]

	err := temporal.SignalWorkflow(r.Context(), workflowID, "", signalName, update)

	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	res := models.ErrorResponse{Message: "Endpoint not found"}
	json.NewEncoder(w).Encode(res)
}

func WriteBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	res := models.ErrorResponse{Message: err.Error()}
	json.NewEncoder(w).Encode(res)
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	res := models.ErrorResponse{Message: err.Error()}
	json.NewEncoder(w).Encode(res)
}
