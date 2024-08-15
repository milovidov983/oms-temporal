PUT http://localhost:8888/order/ORDER-1723715599/assembly
Content-Type: application/json
[
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

HTTP/1.1 200
