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
		CancelReason    string
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
	OrderStatusPassedToAssembly
	OrderStatusAssembled
	OrderStatusDelivering
	OrderStatusDelivered
	OrderStatusCanceled
)

var statusName = map[OrderStatus]string{
	OrderStatusCreated:          "created",
	OrderStatusPassedToAssembly: "passed_to_assembly",
	OrderStatusAssembled:        "assembled",
	OrderStatusDelivering:       "delivering",
	OrderStatusDelivered:        "delivered",
	OrderStatusCanceled:         "canceled",
}

func (os OrderStatus) String() string {
	return statusName[os]
}
