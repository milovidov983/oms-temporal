package main

import (
	"log"

	"github.com/milovidov983/oms-temporal/internal/cartorder"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: "127.0.0.1:7777",
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, "ORDER_TASK_QUEUE", worker.Options{})

	a := &cartorder.Activities{}

	w.RegisterActivity(a.UpdateDeliveryComment)
	w.RegisterActivity(a.UpdateAssemblyComment)
	w.RegisterActivity(a.SendOrderToAssembly)
	w.RegisterActivity(a.SendOrderToDelivery)
	w.RegisterActivity(a.SendEventOrderStatusChanged)
	w.RegisterActivity(a.SendEventAssemblyCommentChanged)
	w.RegisterActivity(a.SendEventAssemblyCommentFailedToChange)
	w.RegisterActivity(a.SendEventDeliveryCommentChanged)
	w.RegisterActivity(a.SendEventDeliveryCommentFailedToChange)

	w.RegisterWorkflow(cartorder.CartOrderWorkflow)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
