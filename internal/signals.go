package internal

import "github.com/milovidov983/oms-temporal/pkg/models"

// var SignalChannels = struct {
// 	COMPLETE_ASSEMBLY_CHANNEL       string
// 	CHANGE_ASSEMBLY_COMMENT_CHANNEL string
// 	COMPLETE_DELIVERY_CHANNEL       string
// 	CHANGE_DELIVERY_COMMENT_CHANNEL string
// 	CANCEL_ORDER_CHANNEL            string
// }{
// 	COMPLETE_ASSEMBLY_CHANNEL:       "COMPLETE_ASSEMBLY_CHANNEL",
// 	CHANGE_ASSEMBLY_COMMENT_CHANNEL: "CHANGE_ASSEMBLY_COMMENT_CHANNEL",
// 	COMPLETE_DELIVERY_CHANNEL:       "COMPLETE_DELIVERY_CHANNEL",
// 	CHANGE_DELIVERY_COMMENT_CHANNEL: "CHANGE_DELIVERY_COMMENT_CHANNEL",
// 	CANCEL_ORDER_CHANNEL:            "CANCEL_ORDER_CHANNEL",
// }

var RouteTypes = struct {
	COMPLETE_ASSEMBLY       string
	CHANGE_ASSEMBLY_COMMENT string
	COMPLETE_DELIVERY       string
	CHANGE_DELIVERY_COMMENT string
	CANCEL_ORDER            string
}{
	COMPLETE_ASSEMBLY:       "complete_assembly",
	CHANGE_ASSEMBLY_COMMENT: "change_assembly_comment",
	COMPLETE_DELIVERY:       "complete_delivery",
	CHANGE_DELIVERY_COMMENT: "change_delivery_comment",
	CANCEL_ORDER:            "cancel_order",
}

type RouteSignal struct {
	Route string
}

type SignalPayloadCompleteAssembly struct {
	Route     string
	Collected models.OrderLines
}

type SignalPayloadChangeAssemblyComment struct {
	Route   string
	Comment string
}

type SignalPayloadCompleteDelivery struct {
	Route     string
	Delivered models.OrderLines
}

type SignalPayloadChangeDeliveryComment struct {
	Route   string
	Comment string
}

type SignalPayloadCancelOrder struct {
	Route  string
	Reason string
}
