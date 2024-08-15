POST http://localhost:8888/order
Content-Type: application/json
{
  "assembly_comment": "Будьте осторожны с хрупкими предметами",
  "delivery_comment": "Позвоните за час до доставки",
  "ordered": [
    {
      "product_id": 1,
      "quantity": 2,
      "price": 99.99
    },
    {
      "product_id": 3,
      "quantity": 1,
      "price": 49.99
    }
  ]
}

HTTP/1.1 201
[Asserts]
jsonpath "$.order" exists
jsonpath "$.order.order_id" exists
jsonpath "$.order.assembly_comment" == "Будьте осторожны с хрупкими предметами"
jsonpath "$.order.delivery_comment" == "Позвоните за час до доставки"
jsonpath "$.order.ordered" count == 2
jsonpath "$.workflowID" exists
jsonpath "$.workflowID" matches "ORDER-[0-9]+"