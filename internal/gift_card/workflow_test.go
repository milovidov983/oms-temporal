package gift_card

import (
	"errors"
	"testing"

	"github.com/milovidov983/oms-temporal/external"
	"github.com/milovidov983/oms-temporal/pkg/models"
	"github.com/milovidov983/oms-temporal/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessProcessingWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("", nil)
	env.OnActivity(GetGiftCardNumber, mock.Anything, testInput).Return("", nil)
	env.OnActivity(SendSuccessNotification, mock.Anything, testInput).Return(nil)
	env.OnActivity(ExecuteWebsiteCallback, mock.Anything, testInput, "").Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_PayFailed_InsufficientFundsError_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("", &external.InsufficientFundsError{})
	env.OnActivity(SendFailureNotification, mock.Anything, testInput).Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_FailedGiftCardCreation_SuccessfulRefund_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("transaction-id", nil)
	env.OnActivity(GetGiftCardNumber, mock.Anything, testInput).Return("", errors.New("Не смогли создать создать подарочную карту"))
	env.OnActivity(Refund, mock.Anything, "transaction-id", testInput.Metadata.IdempotencyToken).Return("", nil)
	env.OnActivity(SendRefundNotification, mock.Anything, testInput).Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_FailedGiftCardCreation_FailedRefund_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("transaction-id", nil)
	env.OnActivity(GetGiftCardNumber, mock.Anything, testInput).Return("", errors.New("Не смогли создать создать подарочную карту"))
	env.OnActivity(Refund, mock.Anything, "transaction-id", testInput.Metadata.IdempotencyToken).Return("", errors.New("Не смогли вернуть деньги"))
	env.OnActivity(SendSupportAlert, mock.Anything, testInput, "Failed to refund transaction! Please solve this problem!").Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_FailedCallback_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("transaction-id", nil)
	env.OnActivity(GetGiftCardNumber, mock.Anything, testInput).Return("gift-card-number", nil)
	env.OnActivity(SendSuccessNotification, mock.Anything, testInput).Return(nil)
	env.OnActivity(ExecuteWebsiteCallback, mock.Anything, testInput, "gift-card-number").Return(errors.New("Сайт умер"))
	env.OnActivity(SendSupportAlert, mock.Anything, testInput, "Failed to return a gift card number using a callback!").Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_FailedNotification_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testInput := createTestInput()

	env.OnActivity(Pay, mock.Anything, testInput).Return("transaction-id", nil)
	env.OnActivity(GetGiftCardNumber, mock.Anything, testInput).Return("gift-card-number", nil)
	env.OnActivity(SendSuccessNotification, mock.Anything, testInput).Return(errors.New("Пуши умерли"))
	env.OnActivity(ExecuteWebsiteCallback, mock.Anything, testInput, "gift-card-number").Return(nil)

	env.ExecuteWorkflow(Processing, testInput)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func createTestInput() models.GiftCardOrderRequest {
	return models.GiftCardOrderRequest{
		OrderID:     10043,
		CardType:    "GIFT",
		Amount:      250,
		CallbackURL: "https://microsoft.com/xxx",
		Customer: models.Customer{
			CustomerID: "999",
		},
		Payment: models.PaymentDetails{
			AccountNumber: "11-111",
			Amount:        250,
		},
		Metadata: models.Metadata{
			IdempotencyToken: utils.Pseudo_uuid(),
		},
	}
}
