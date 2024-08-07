package external

import "math/rand"

type GiftCardService struct {
}

type GiftCardRequest struct {
	OrderID int
}

func (s *GiftCardService) CreateGiftCard(orderID int, idempotencyToken string) (string, error) {
	giftCardNumber := generateGiftCardNumber()
	return giftCardNumber, nil
}

func generateGiftCardNumber() string {
	randChars := make([]byte, 10)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return "GIFT_CARD_NUMBER_" + string(randChars)
}
