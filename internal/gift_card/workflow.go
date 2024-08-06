package gift_card

import (
	"time"

	"github.com/milovidov983/oms-temporal/pkg/models"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func IssueGiftCard(ctx workflow.Context, input models.GiftCardOrderRequest) (string, error) {

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        42, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	var paymentOutput string

	paymentErr := workflow.ExecuteActivity(ctx, Pay, input).Get(ctx, &paymentOutput)

	if paymentErr != nil {
		notificationErr := workflow.ExecuteActivity(ctx, SendFailureNotification, input).Get(ctx, nil)
		if notificationErr != nil {
			return "", notificationErr
		}
		return "", paymentErr
	}
	notificationErr := workflow.ExecuteActivity(ctx, SendSuccessNotification, input).Get(ctx, nil)
	if notificationErr != nil {
		return "", notificationErr
	}

	return paymentOutput, nil
}
