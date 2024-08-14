package giftcardorder

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
	transactionID, err := payService.MakePayment(data.Payment.AccountNumber, data.Amount, it)
	log.Printf("Pay service API called. Payment confirmation: %s\n\n", transactionID)
	return transactionID, err
}

func Refund(ctx context.Context, transactionID string, idempotencyToken string) (string, error) {
	log.Printf("Refunding gift card. TransactionID %s.\n\n", transactionID)

	payService := external.PayService{}
	transactionID, err := payService.RefundPayment(transactionID, idempotencyToken)

	log.Printf("Refund service API called.\n\n")

	return transactionID, err
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

func SendSupportAlert(ctx context.Context, data models.GiftCardOrderRequest, text string) error {
	log.Printf("SUPPORT NOTIFICATION: %s! CustomerID %s\n\n", text, data.Customer.CustomerID)
	return nil
}

func SendRefundNotification(ctx context.Context, data models.GiftCardOrderRequest) error {
	log.Printf("Sending failure notification to %s.\n\n", data.Customer.CustomerID)

	text := "Your gift card purchase failed! Your funds will be refunded to your card"
	notificationService := external.NotificationService{}
	err := notificationService.Send(&external.PushRequest{
		To:   data.Customer.CustomerID,
		Text: text,
	}, data.Metadata.IdempotencyToken)

	log.Printf("Notification service API called.\n\n")
	return err
}

func ExecuteWebsiteCallback(ctx context.Context, data models.GiftCardOrderRequest, giftCardNumber string) error {
	log.Printf("Executing website callback for order %d. Gift card number %s\n\n", data.OrderID, giftCardNumber)

	callbackService := CallbackService{CallbackURL: data.CallbackURL}
	err := callbackService.Execute(giftCardNumber, data.Metadata.IdempotencyToken)

	log.Printf("Callback service API called.\n\n")
	return err
}

func GetGiftCardNumber(ctx context.Context, data models.GiftCardOrderRequest) (string, error) {
	log.Printf("Creating gift card number for order %d.\n\n", data.OrderID)

	giftCardService := external.GiftCardService{}
	giftCardNumber, err := giftCardService.CreateGiftCard(data.OrderID, data.Metadata.IdempotencyToken)

	log.Printf("Gift card service API called. Gift card number: %s\n\n", giftCardNumber)
	return giftCardNumber, err
}
