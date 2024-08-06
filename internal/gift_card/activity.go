package gift_card

import (
	"context"
	"fmt"
	"log"

	"github.com/milovidov983/oms-temporal/external"
	"github.com/milovidov983/oms-temporal/pkg/models"
)

func Pay(ctx context.Context, data models.GiftCardOrderRequest) (string, error) {
	log.Printf("Paying $%d from account %s.\n\n",
		data.Amount,
		data.Payment.AccountNumber,
	)

	it := fmt.Sprintf("%s-payment", data.Metadata.IdempotencyToken)
	payService := external.PayService{}
	confirmation, err := payService.MakePayment(data.Payment.AccountNumber, data.Amount, it)
	log.Printf("Pay service API called. Payment confirmation: %s\n\n", confirmation)
	return confirmation, err
}

func SendSuccessNotification(ctx context.Context, data models.GiftCardOrderRequest) error {
	log.Printf("Sending success notification to %s.\n\n", data.Customer.CustomerID)

	text := "Your gift card purchase was successful!"
	notificationService := external.NotificationService{}
	err := notificationService.Send(&external.PushRequest{
		To:   data.Customer.CustomerID,
		Text: text,
	}, data.Metadata.IdempotencyToken)

	log.Printf("Notification service API called.\n\n")
	return err
}

func SendFailureNotification(ctx context.Context, data models.GiftCardOrderRequest) error {
	log.Printf("Sending failure notification to %s.\n\n", data.Customer.CustomerID)

	text := "Your gift card purchase failed!"
	notificationService := external.NotificationService{}
	err := notificationService.Send(&external.PushRequest{
		To:   data.Customer.CustomerID,
		Text: text,
	}, data.Metadata.IdempotencyToken)

	log.Printf("Notification service API called.\n\n")
	return err
}
