package routes

const RouteTypeCreateOrder = "create_order"
const RouteTypeCompleteAssembly = "complete_assembly"
const RouteTypeChangeAssemblyComment = "change_assembly_comment"
const RouteTypeCompleteDelivery = "complete_delivery"
const RouteTypeChangeDeliveryComment = "change_delivery_comment"
const RouteTypeCancelOrder = "cancel_order"

// Не реализовано:
// Надо сделать child workflow для процессов передачи
// заказа в сборку и в доставку и всю логику реализовать в
// отдельных child workflow
const RouteTypeTransferToAssembly = "transfer_to_assembly"
const RouteTypeTransferToDelivery = "transfer_to_delivery"
