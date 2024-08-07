package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/milovidov983/oms-temporal/pkg/models"
)

// @@@SNIPSTART money-transfer-project-template-go-start-workflow
func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := models.GiftCardOrderRequest{
		OrderID:     10042,
		CardType:    "GIFT",
		Amount:      250,
		CallbackURL: "https://microsoft.com/xxx",
		Customer: models.Customer{
			CustomerID: "99",
		},
		Payment: models.PaymentDetails{
			AccountNumber: "42000042",
			Amount:        250,
		},
		Metadata: models.Metadata{
			IdempotencyToken: "",
		},
	}

	options := client.StartWorkflowOptions{
		ID:        "gift-card-order-042",
		TaskQueue: models.GiftCardTaskQueueName,
	}

	log.Printf("Starting order gift card workflow from account %s to account %s for %d", input.SourceAccount, input.TargetAccount, input.Amount)

	we, err := c.ExecuteWorkflow(context.Background(), options, app.MoneyTransfer, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}

// @@@SNIPEND
