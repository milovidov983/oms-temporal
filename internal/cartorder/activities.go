package cartorder

import (
	"context"
	"log"

	"github.com/milovidov983/oms-temporal/pkg/models"
)

type Activities struct{}

func (a *Activities) UpdateDeliveryComment(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Delivery comment updated for order %s", order.OrderID)
	return nil
}

func (a *Activities) UpdateAssemblyComment(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Assembly comment updated for order %s", order.OrderID)
	return nil
}

func (a *Activities) SendOrderToAssembly(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Order %s sent to assembly", order.OrderID)
	return nil
}

func (a *Activities) SendOrderToDelivery(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Order %s sent to delivery", order.OrderID)
	return nil
}

func (a *Activities) SendEventOrderStatusChanged(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Event order status changed for order %s new status is %s", order.OrderID, order.Status.String())
	return nil
}

func (a *Activities) SendEventAssemblyCommentChanged(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Event assembly comment changed for order %s", order.OrderID)
	return nil
}

func (a *Activities) SendEventAssemblyCommentFailedToChange(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Event assembly comment failed to change for order %s, current status is %s", order.OrderID, order.Status.String())
	return nil
}

func (a *Activities) SendEventDeliveryCommentChanged(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Event delivery comment changed for order %s", order.OrderID)
	return nil
}

func (a *Activities) SendEventDeliveryCommentFailedToChange(ctx context.Context, order models.OrderState) error {
	log.Printf("[info] Event delivery comment failed to change for order %s, current status is %s", order.OrderID, order.Status.String())
	return nil
}
