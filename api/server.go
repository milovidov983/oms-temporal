package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bojanz/httpx"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/milovidov983/oms-temporal/internal"
	"github.com/milovidov983/oms-temporal/internal/cartorder"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"go.temporal.io/sdk/client"
)

type (
	ErrorResponse struct {
		Message string
	}
)

var (
	HTTPPort = os.Getenv("PORT")
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
	r.Handle("/order/{workflowID}/assembly-comment", http.HandlerFunc(AssemblyCommentOrderHandler)).Methods("PUT")
	r.Handle("/cart/{workflowID}/delivery", http.HandlerFunc(DeliveryOrderHandler)).Methods("PUT")
	r.Handle("/cart/{workflowID}/delivery-comment", http.HandlerFunc(DeliveryCommentOrderHandler)).Methods("PUT")
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
		OrderID:   orderID,
		Status:    models.OrderStatusCreated,
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
	var collected models.OrderLines
	if err := json.NewDecoder(r.Body).Decode(&collected); err != nil {
		WriteError(w, err)
		return
	}

	update := internal.SignalPayloadCompleteAssembly{
		Route:     internal.RouteTypes.COMPLETE_ASSEMBLY,
		Collected: collected,
	}

	err := temporal.SignalWorkflow(r.Context(), vars["workflowID"], "", internal.SignalChannels.COMPLETE_ASSEMBLY_CHANNEL, update)
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	res := make(map[string]interface{})
	res["ok"] = 1
	json.NewEncoder(w).Encode(res)
}

func ChangeAssemblyComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var comment string
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		WriteError(w, err)
		return
	}

	update := internal.SignalPayloadChangeAssemblyComment{
		Route:   internal.RouteTypes.CHANGE_ASSEMBLY_COMMENT,
		Comment: comment,
	}

	signalName := internal.SignalChannels.CHANGE_ASSEMBLY_COMMENT_CHANNEL
	err := temporal.SignalWorkflow(r.Context(), vars["workflowID"], "", signalName, update)

	if errors.Is(err, internal.ErrWrongStatus) {
		WriteBadRequest(w, err)
		return
	}
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
	res := ErrorResponse{Message: "Endpoint not found"}
	json.NewEncoder(w).Encode(res)
}

func WriteBadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	res := ErrorResponse{Message: err.Error()}
	json.NewEncoder(w).Encode(res)
}

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	res := ErrorResponse{Message: err.Error()}
	json.NewEncoder(w).Encode(res)
}
