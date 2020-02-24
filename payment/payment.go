package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bosamatheus/ecommerce-microservices/payment/queue"
	amqp "github.com/streadway/amqp"
)

type Order struct {
	UUID      string    `json:"uuid"`
	ProductID string    `json:"product_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	in := make(chan []byte)

	conn := queue.Connect()
	queue.StartConsuming("order_queue", conn, in)

	var order Order
	// Draining out the channel in
	for payload := range in {
		json.Unmarshal(payload, &order)
		order.Status = "approved"
		notifyPaymentProcessed(order, conn)
	}
}

func notifyPaymentProcessed(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, "payment_ex", "", ch)
	fmt.Println(string(json))
}
