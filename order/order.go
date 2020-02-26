package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bosamatheus/ecommerce-microservices/order/db"
	"github.com/bosamatheus/ecommerce-microservices/order/queue"
	uuid "github.com/nu7hatch/gouuid"
	amqp "github.com/streadway/amqp"
)

type Product struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"name"`
	Price float64 `json:"price,string"`
}

type Order struct {
	UUID      string    `json:"uuid"`
	ProductID string    `json:"product_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var productURL string

// Executes before main
func init() {
	// productURL := os.Getenv("PRODUCT_URL")
	productURL = "http://localhost:8082"
	_ = productURL
}

func main() {
	var param string
	flag.StringVar(&param, "opt", "", "Usage")
	flag.Parse()

	in := make(chan []byte)

	conn := queue.Connect()

	switch param {
	case "checkout":
		queue.StartConsuming("checkout_queue", conn, in)

		// Draining out the channel
		for payload := range in {
			order := createOrder(payload)
			notifyOrderCreated(order, conn)

			fmt.Println("Checkout: ", string(payload))
		}
	case "payment":
		queue.StartConsuming("payment_queue", conn, in)

		var order Order
		// Draining out the channel
		for payload := range in {
			json.Unmarshal(payload, &order)
			saveOrder(order)

			fmt.Println("Payment: ", string(payload))
		}
	}
}

func getProductByID(id string) Product {
	response, err := http.Get(productURL + "/product/" + id)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		fmt.Printf("An error has occurred while unmarshal json object: %s\n", err)
	}

	return product
}

func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.UUID = uuid.String()
	order.Status = "pendent"
	order.CreatedAt = time.Now()
	saveOrder(order)
	return order
}

func saveOrder(order Order) {
	json, _ := json.Marshal(order)
	conn := db.Connect()

	err := conn.Set(order.UUID, string(json), 0).Err()
	if err != nil {
		panic(err.Error())
	}
}

func notifyOrderCreated(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	queue.Notify(json, "order_ex", "", ch)
}
