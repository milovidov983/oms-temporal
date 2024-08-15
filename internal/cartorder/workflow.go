package cartorder

import (
	"time"

	"github.com/milovidov983/oms-temporal/internal/signals"
	"github.com/milovidov983/oms-temporal/internal/signals/channels"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CartOrderWorkflow(ctx workflow.Context, state models.OrderState) error {
	// Настраиваем политику ретраев
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}
	// настраиваем опции запуска всех активити
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	logger := workflow.GetLogger(ctx)

	// Обработка запросов в workflow
	// Запрос на получение заказа
	err := workflow.SetQueryHandler(ctx, "getOrder", func(input []byte) (models.OrderState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	// Создание каналов
	createOrderChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCreateOrderChannel)
	transferToAssemblyLocalChannel := workflow.NewChannel(ctx)
	completeAssemblyChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCompleteAssemblyChannel)
	transferToDeliveryLocalChannel := workflow.NewChannel(ctx)
	changeAssemblyCommentChannel := workflow.GetSignalChannel(ctx, channels.SignalNameChangeAssemblyCommentChannel)
	completeDeliveryChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCompleteDeliveryChannel)
	changeDeliveryCommentChannel := workflow.GetSignalChannel(ctx, channels.SignalNameChangeDeliveryCommentChannel)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, channels.SignalNameCancelOrderChannel)

	orderCompleted := false

	a := &Activities{}

	for {
		selector := workflow.NewSelector(ctx)
		// Заказ создан
		selector.AddReceive(createOrderChannel, func(c workflow.ReceiveChannel, _ bool) {

			state.Status = models.OrderStatusCreated

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}

			transferToAssemblyLocalChannel.Send(ctx, struct{}{})
		})
		// Отправка в сборку
		selector.AddReceive(transferToAssemblyLocalChannel, func(c workflow.ReceiveChannel, _ bool) {

			sendAssemblyErr := workflow.ExecuteActivity(ctx, a.SendOrderToAssembly, state).Get(ctx, nil)
			if sendAssemblyErr != nil {
				logger.Error("Failed to send order to assembly: %v", sendAssemblyErr)
				return
			}

			state.Status = models.OrderStatusTransferredToAssembly

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}
		})
		// Завершение сборки
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

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}

			// Отправили заказ в доставку
			transferToDeliveryLocalChannel.Send(ctx, struct{}{})
		})
		// Отправка в доставку
		selector.AddReceive(transferToDeliveryLocalChannel, func(c workflow.ReceiveChannel, _ bool) {

			sendToDeliveryErr := workflow.ExecuteActivity(ctx, a.SendOrderToDelivery, state).Get(ctx, nil)
			if sendToDeliveryErr != nil {
				logger.Error("Failed to send order to delivery: %v", sendToDeliveryErr)
				return
			}

			state.Status = models.OrderStatusTransferredToDelivery

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}
		})
		// Изменение комментария для сборщика
		selector.AddReceive(changeAssemblyCommentChannel, func(c workflow.ReceiveChannel, _ bool) {
			isValidStatus := state.Status.Any(models.OrderStatusCreated, models.OrderStatusTransferredToAssembly)
			if !isValidStatus {
				sendEventCommentErr := workflow.ExecuteActivity(ctx, a.SendEventAssemblyCommentFailedToChange, state)
				if sendEventCommentErr != nil {
					logger.Error("Failed to send assembly comment failed to change event: %v", sendEventCommentErr)
				}
				logger.Warn("Invalid status for changing assembly comment: %v", state.Status)
				return
			}

			var signal interface{}
			c.Receive(ctx, &signal)

			var message signals.SignalPayloadChangeAssemblyComment
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			state.AssemblyComment = message.Comment

			sendEventCommentErr := workflow.ExecuteActivity(ctx, a.SendEventAssemblyCommentChanged, state).Get(ctx, nil)
			if sendEventCommentErr != nil {
				logger.Error("Failed to send assembly comment changed event: %v", sendEventCommentErr)
			}
		})
		// Изменение комментария для доставщика
		selector.AddReceive(changeDeliveryCommentChannel, func(c workflow.ReceiveChannel, _ bool) {
			validStatuses := map[models.OrderStatus]bool{
				models.OrderStatusCreated:               true,
				models.OrderStatusTransferredToAssembly: true,
				models.OrderStatusAssemblyInProgress:    true,
				models.OrderStatusAssembled:             true,
				models.OrderStatusTransferredToDelivery: true,
			}
			isValidStatus := validStatuses[state.Status]
			if !isValidStatus {
				sendEventCommentErr := workflow.ExecuteActivity(ctx, a.SendEventDeliveryCommentFailedToChange, state)
				if sendEventCommentErr != nil {
					logger.Error("Failed to send delivery comment failed to change event: %v", sendEventCommentErr)
				}
				logger.Warn("Invalid status for changing delivery comment: %v", state.Status)
				return
			}

			var signal interface{}
			c.Receive(ctx, &signal)

			var message signals.SignalPayloadChangeDeliveryComment
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			state.DeliveryComment = message.Comment

			sendEventCommentErr := workflow.ExecuteActivity(ctx, a.SendEventDeliveryCommentChanged, state).Get(ctx, nil)
			if sendEventCommentErr != nil {
				logger.Error("Failed to send delivery comment changed event: %v", sendEventCommentErr)
			}
		})
		// Завершение доставки
		selector.AddReceive(completeDeliveryChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message signals.SignalPayloadCompleteDelivery
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			state.Status = models.OrderStatusDelivered

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}

			orderCompleted = true
		})

		// Отмена заказа
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

			sendEventStatusErr := workflow.ExecuteActivity(ctx, a.SendEventOrderStatusChanged, state).Get(ctx, nil)
			if sendEventStatusErr != nil {
				logger.Error("Failed to send order status changed event: %v", sendEventStatusErr)
			}

			orderCompleted = true
		})

		selector.Select(ctx)

		if orderCompleted {
			break
		}
	}

	return nil
}
