package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Product struct {
	UUID    string  `json: "uuid"`
	Product string  `json: "product"`
	Price   float64 `json: "price, string"`
}

type Products struct {
	Products []Product
}

func loadData() []byte {
	jsonFile, err := os.Open("products.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)

	return data
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	products := loadData()
	w.Write([]byte(products))
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := loadData()

	var products Products
	err := json.Unmarshal(data, &products)

	if err == nil {
		for _, v := range products.Products {
			if v.UUID == vars["id"] {
				product, _ := json.Marshal(v)
				w.Write([]byte(product))
			}
		}
	} else {
		fmt.Println(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", ListProducts)
	r.HandleFunc("/products/{id}", GetProductById)

	http.ListenAndServe(":8081", r)
}
