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
	UUID  string  `json: "uuid"`
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

type Products struct {
	Products []Product `json: "products"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/products", listProducts)
	router.HandleFunc("/products/{id}", getProductById)

	http.ListenAndServe(":8082", router)
}

func loadData() []byte {
	jsonFile, err := os.Open("products.json")
	if err != nil {
		fmt.Printf("An error has occurred while loading json file: %s\n", err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)

	return data
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	products := loadData()
	w.Write([]byte(products))
}

func getProductById(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("An error has occurred while unmarshal json object: %s\n", err)
	}
}
