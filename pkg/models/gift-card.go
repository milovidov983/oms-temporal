package models

const GiftCardTaskQueueName = "GIFT_CARD_TASK_QUEUE"

type GiftCardOrderRequest struct {
	OrderID  int
	CardType string
	Amount   int
	Customer Customer
	Payment  PaymentDetails
	Metadata Metadata
}

type Customer struct {
	CustomerID string
}

type PaymentDetails struct {
	AccountNumber string
	Amount        int
}

type Metadata struct {
	IdempotencyToken string
}
