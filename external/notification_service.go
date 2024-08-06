package external

type NotificationService struct {
}

type PushRequest struct {
	To   string
	Text string
}

func (s *NotificationService) Send(req *PushRequest, idempotencyToken string) error {

	return nil
}
