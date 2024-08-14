package models

type (
	OrderState struct {
		OrderID         string
		Status          OrderStatus
		AssemblyComment string
		DeliveryComment string
		Ordered         []OrderLines
		Collected       []OrderLines
		Delivered       []OrderLines
	}
	OrderLines struct {
		ProductID int
		Quantity  int
		Price     float64
	}
)

type OrderStatus int

const (
	OrderStatusCreated = iota
	OrderStatusAssembling
	OrderStatusAssembled
	OrderStatusDelivering
	OrderStatusDelivered
)

var statusName = map[OrderStatus]string{
	OrderStatusCreated:    "created",
	OrderStatusAssembling: "assembling",
	OrderStatusAssembled:  "assembled",
	OrderStatusDelivering: "delivering",
	OrderStatusDelivered:  "delivered",
}

func (os OrderStatus) String() string {
	return statusName[os]
}
