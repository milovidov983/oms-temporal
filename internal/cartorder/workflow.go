package cartorder

import (
	"github.com/milovidov983/oms-temporal/internal/signals"
	"github.com/milovidov983/oms-temporal/internal/signals/channels"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

func CartOrderWorkflow(ctx workflow.Context, state models.OrderState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "getOrder", func(input []byte) (models.OrderState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	completeAssemblyChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCompleteAssemblyChannel)
	// changeAssemblyCommentChannel := workflow.GetSignalChannel(ctx, channels.SignalNameChangeAssemblyCommentChannel)
	// completeDeliveryChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCompleteDeliveryChannel)
	// changeDeliveryCommentChannel := workflow.GetSignalChannel(ctx, channels.SignalNameChangeDeliveryCommentChannel)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCancelOrderChannel)

	orderCompleted := false

	a := &Activities{
		Logger: logger,
	}

	// Отправили заказа в сборку
	a.SendOrderToAssembly(state)
	state.Status = models.OrderStatusPassedToAssembly
	a.SendEventOrderStatusChanged(state)

	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(completeAssemblyChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message signals.SignalPayloadCompleteAssembly
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			state.Collected = message.Collected
			state.Status = models.OrderStatusAssembled
			a.SendEventOrderStatusChanged(state)
			a.SendOrderToDelivery(state)
		})
		selector.AddReceive(cancelOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message signals.SignalPayloadCancelOrder
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			state.Status = models.OrderStatusCanceled
			state.CancelReason = message.Reason
			a.SendEventOrderStatusChanged(state)
			orderCompleted = true
		})

		selector.Select(ctx)

		if orderCompleted {
			break
		}
	}

	return nil
}
