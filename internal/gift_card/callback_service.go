package gift_card

type CallbackService struct {
	CallbackURL string
}

func (s *CallbackService) Execute(params string, idempotencyToken string) error {
	// send data to the callback resource
	return nil
}
