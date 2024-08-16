package models

type (
	OrderState struct {
		OrderID         string       `json:"order_id"`
		Status          OrderStatus  `json:"status"`
		AssemblyComment string       `json:"assembly_comment"`
		DeliveryComment string       `json:"delivery_comment"`
		Ordered         []OrderLines `json:"ordered"`
		Collected       []OrderLines `json:"collected"`
		Delivered       []OrderLines `json:"delivered"`
		CancelReason    string       `json:"cancel_reason"`
	}
	OrderLines struct {
		ProductID int     `json:"product_id"`
		Quantity  int     `json:"quantity"`
		Price     float64 `json:"price"`
	}
)

type OrderStatus int

const (
	OrderStatusUndefined = iota
	OrderStatusCreated
	OrderStatusTransferredToAssembly
	OrderStatusAssemblyInProgress
	OrderStatusAssembled
	OrderStatusTransferredToDelivery
	OrderStatusDeliveryInProgress
	OrderStatusDelivered
	OrderStatusCanceled
)

var statusName = map[OrderStatus]string{
	OrderStatusUndefined:             "undefined",
	OrderStatusCreated:               "created",
	OrderStatusTransferredToAssembly: "transferred_to_assembly",
	OrderStatusAssemblyInProgress:    "assembly_in_progress",
	OrderStatusAssembled:             "assembled",
	OrderStatusTransferredToDelivery: "transferred_to_delivery",
	OrderStatusDeliveryInProgress:    "delivery_in_progress",
	OrderStatusDelivered:             "delivered",
	OrderStatusCanceled:              "canceled",
}

func (os OrderStatus) String() string {
	return statusName[os]
}

func (os OrderStatus) Any(args ...OrderStatus) bool {
	for _, a := range args {
		if os == a {
			return true
		}
	}
	return false
}
