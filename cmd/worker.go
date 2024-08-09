package main

import (
	"log"

	"github.com/milovidov983/oms-temporal/internal/gift_card"
	"github.com/milovidov983/oms-temporal/pkg/models"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, models.GiftCardTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(gift_card.Processing)
	w.RegisterActivity(gift_card.Pay)
	w.RegisterActivity(gift_card.Refund)
	w.RegisterActivity(gift_card.SendSuccessNotification)
	w.RegisterActivity(gift_card.SendFailureNotification)
	w.RegisterActivity(gift_card.SendSupportAlert)
	w.RegisterActivity(gift_card.SendRefundNotification)
	w.RegisterActivity(gift_card.ExecuteWebsiteCallback)
	w.RegisterActivity(gift_card.GetGiftCardNumber)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
