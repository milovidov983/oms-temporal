package cartorder

import (
	"go.temporal.io/sdk/log"

	"github.com/milovidov983/oms-temporal/pkg/models"
)

type Activities struct {
	Logger log.Logger
}

func (a *Activities) UpdateDeliveryComment(order models.OrderState) error {
	a.Logger.Info("Delivery comment updated for order %s", order.OrderID)
	return nil
}

func (a *Activities) UpdateAssemblyComment(order models.OrderState) error {
	a.Logger.Info("Assembly comment updated for order %s", order.OrderID)
	return nil
}

func (a *Activities) SendOrderToAssembly(order models.OrderState) error {
	a.Logger.Info("Order %s sent to assembly", order.OrderID)
	return nil
}

func (a *Activities) SendOrderToDelivery(order models.OrderState) error {
	a.Logger.Info("Order %s sent to delivery", order.OrderID)
	return nil
}

func (a *Activities) CheckCollectedLines(order models.OrderState) error {
	a.Logger.Info("Collected lines checked for order %s", order.OrderID)
	return nil
}

func (a *Activities) SendEventOrderStatusChanged(order models.OrderState) error {
	a.Logger.Info("Event order status changed for order %s new status is %s", order.OrderID, order.Status.String())
	return nil
}
