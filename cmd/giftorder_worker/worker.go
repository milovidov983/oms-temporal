package main

import (
	"log"

	"github.com/milovidov983/oms-temporal/internal/giftcardorder"
	"github.com/milovidov983/oms-temporal/pkg/models"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, models.GiftCardTaskQueueName, worker.Options{})

	w.RegisterWorkflow(giftcardorder.Processing)
	w.RegisterActivity(giftcardorder.Pay)
	w.RegisterActivity(giftcardorder.Refund)
	w.RegisterActivity(giftcardorder.SendSuccessNotification)
	w.RegisterActivity(giftcardorder.SendFailureNotification)
	w.RegisterActivity(giftcardorder.SendSupportAlert)
	w.RegisterActivity(giftcardorder.SendRefundNotification)
	w.RegisterActivity(giftcardorder.ExecuteWebsiteCallback)
	w.RegisterActivity(giftcardorder.GetGiftCardNumber)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
