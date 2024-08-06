package gift_card

type PaymentDetails struct {
	AccountNumber string
	Amount        int
	Metadata      Metadata
}

type Metadata struct {
	IdempotencyToken string
}
