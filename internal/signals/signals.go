package signals

import "github.com/milovidov983/oms-temporal/pkg/models"

type RouteSignal struct {
	Route string
}

type SignalPayloadCompleteAssembly struct {
	Route     string
	Collected []models.OrderLines
}

type SignalPayloadChangeAssemblyComment struct {
	Route   string
	Comment string
}

type SignalPayloadCompleteDelivery struct {
	Route     string
	Delivered []models.OrderLines
}

type SignalPayloadChangeDeliveryComment struct {
	Route   string
	Comment string
}

type SignalPayloadCancelOrder struct {
	Route  string
	Reason string
}
