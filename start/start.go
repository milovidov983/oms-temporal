package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/milovidov983/oms-temporal/internal/giftcardorder"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"github.com/milovidov983/oms-temporal/pkg/utils"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: "127.0.0.1:7777",
	})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := models.GiftCardOrderRequest{
		OrderID:     10043,
		CardType:    "GIFT",
		Amount:      250,
		CallbackURL: "https://microsoft.com/xxx",
		Customer: models.Customer{
			CustomerID: "999",
		},
		Payment: models.PaymentDetails{
			AccountNumber: "11-111",
			Amount:        250,
		},
		Metadata: models.Metadata{
			IdempotencyToken: utils.Pseudo_uuid(),
		},
	}

	options := client.StartWorkflowOptions{
		ID:        "gift-card-order-045",
		TaskQueue: models.GiftCardTaskQueueName,
	}

	log.Printf("Starting order gift card workflow OrderID %d", input.OrderID)

	we, err := c.ExecuteWorkflow(context.Background(), options, giftcardorder.Processing, input)
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
