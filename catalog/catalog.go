package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type Product struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"name"`
	Price float64 `json:"price,string"`
}

type Products struct {
	Products []Product `json: "products"`
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
	r.HandleFunc("/", listProducts)
	r.HandleFunc("/products/{id}", showProduct)

	http.ListenAndServe(":8083", r)
}

func loadProducts() []Product {
	response, err := http.Get(productsURL + "/products")
	if err != nil {
		fmt.Printf("The HTTP request failed with error: %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	var products Products
	json.Unmarshal(data, &products)

	return products.Products
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := loadProducts()
	t := template.Must(template.ParseFiles("./templates/catalog.html"))
	t.Execute(w, products)
}

func showProduct(w http.ResponseWriter, r *http.Request) {
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

	t := template.Must(template.ParseFiles("./templates/view.html"))
	t.Execute(w, product)
}
