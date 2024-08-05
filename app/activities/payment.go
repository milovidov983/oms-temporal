package activities

import (
	"context"
	"fmt"
	"log"

	"github.com/milovidov983/vvoms/app/models"
)

func Pay(ctx context.Context, data models.PaymentDetails) (string, error) {
	log.Printf("Withdrawing $%d from account %s.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
	bank := BankingService{"bank-api.example.com"}
	confirmation, err := bank.Withdraw(data.SourceAccount, data.Amount, referenceID)
	log.Printf("Bank API called. Withdrawal confirmation: %s\n\n", confirmation)
	return confirmation, err
}
