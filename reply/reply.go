package main

import (
	"context"
	"log"

	"github.com/milovidov983/oms-temporal/internal/giftcardorder"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/history/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()
	err = ReplayWorkflow(context.Background(), c, "d82cc35b-3e15-4e86-8148-9692cb59f030", "6c722bf7-df5b-4e1d-a55b-39a4fa55ec60")
	if err != nil {
		log.Fatalln("Unable to replay workflow:", err)
	}
	log.Println("Workflow replayed successfully")
}

func GetWorkflowHistory(ctx context.Context, client client.Client, id, runID string) (*history.History, error) {
	var hist history.History
	iter := client.GetWorkflowHistory(ctx, id, runID, false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	for iter.HasNext() {
		event, err := iter.Next()
		if err != nil {
			return nil, err
		}
		hist.Events = append(hist.Events, event)
	}
	return &hist, nil
}

func ReplayWorkflow(ctx context.Context, client client.Client, id, runID string) error {
	hist, err := GetWorkflowHistory(ctx, client, id, runID)
	if err != nil {
		return err
	}
	replayer := worker.NewWorkflowReplayer()
	replayer.RegisterWorkflow(giftcardorder.Processing)
	return replayer.ReplayWorkflowHistory(nil, hist)
}
