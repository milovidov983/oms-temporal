package gift_card

import (
	"context"
	"fmt"
	"log"

	"github.com/milovidov983/oms-temporal/external"
)

func Pay(ctx context.Context, data PaymentDetails) (string, error) {
	log.Printf("Paying $%d from account %s.\n\n",
		data.Amount,
		data.AccountNumber,
	)

	it := fmt.Sprintf("%s-make-payment", data.Metadata.IdempotencyToken)
	payService := external.PayService{}
	confirmation, err := payService.MakePayment(data.AccountNumber, data.Amount, it)
	log.Printf("Pay service API called. Payment confirmation: %s\n\n", confirmation)
	return confirmation, err
}
