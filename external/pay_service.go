package external

import (
	"errors"
	"math/rand"
)

type PayService struct {
}

type account struct {
	AccountNumber string
	Balance       int64
}

type bank struct {
	Accounts []account
}

func (b bank) findAccount(accountNumber string) (account, error) {

	for _, v := range b.Accounts {
		if v.AccountNumber == accountNumber {
			return v, nil
		}
	}

	return account{}, errors.New("account not found")
}

type InsufficientFundsError struct{}

func (m *InsufficientFundsError) Error() string {
	return "Insufficient Funds"
}

type InvalidAccountError struct{}

func (m *InvalidAccountError) Error() string {
	return "Account number supplied is invalid"
}

var mockBank = &bank{
	Accounts: []account{
		{AccountNumber: "11-111", Balance: 2000},
		{AccountNumber: "22-222", Balance: 0},
	},
}

func (client PayService) MakePayment(accountNumber string, amount int, idempotencyToken string) (string, error) {
	acct, err := mockBank.findAccount(accountNumber)

	if err != nil {
		return "", &InvalidAccountError{}
	}

	if amount > int(acct.Balance) {
		return "", &InsufficientFundsError{}
	}
	//panic("implement me")
	//return generateTransactionID("P", 10), nil"
	//return "", errors.New("Some err")
	return generateTransactionID("T", 10), nil
}

func (client PayService) RefundPayment(transactionID string, idempotencyToken string) (string, error) {
	// Предполагаем что мы по transactionID понимаем сколько и кому вернуть
	return generateTransactionID("R", 10), nil
}

func generateTransactionID(prefix string, length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
