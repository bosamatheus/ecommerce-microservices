package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/bosamatheus/ecommerce-microservices/checkout/queue"
	"github.com/gorilla/mux"
)

type Product struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"name"`
	Price float64 `json:"price,string"`
}

type Order struct {
	ProductId string `json:"product_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

var productsURL string

// Executes before main
func init() {
	// productsURL := os.Getenv("PRODUCT_URL")
	productsURL = "http://localhost:8082"
	_ = productsURL
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/finish", finish)
	r.HandleFunc("/{id}", displayCheckout)

	http.ListenAndServe(":8084", r)
}

func displayCheckout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response, err := http.Get(productsURL + "/products/" + vars["id"])
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
	}
	data, _ := ioutil.ReadAll(response.Body)

	var product Product
	err = json.Unmarshal(data, &product)
	if err != nil {
		fmt.Printf("An error has occurred while unmarshal json object: %s\n", err)
	}

	t := template.Must(template.ParseFiles("./templates/checkout.html"))
	t.Execute(w, product)
}

func finish(w http.ResponseWriter, r *http.Request) {
	var order Order
	order.ProductId = r.FormValue("product_id")
	order.Name = r.FormValue("name")
	order.Email = r.FormValue("email")
	order.Phone = r.FormValue("phone")

	data, _ := json.Marshal(order)
	fmt.Println(string(data))

	conn := queue.Connect()
	queue.Notify(data, "checkout_ex", "", conn)

	w.Write([]byte("Order processed!"))
}
