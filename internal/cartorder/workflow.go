package cartorder

import (
	"github.com/milovidov983/oms-temporal/pkg/models"
	"go.temporal.io/sdk/workflow"
)

func CartOrderWorkflow(ctx workflow.Context, state models.OrderState) error {
	return nil
}
