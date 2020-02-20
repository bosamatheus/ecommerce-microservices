package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var productsURL string

// Executes before main
func init() {
	// productsURL := os.Getenv("PRODUCT_URL")
	productsURL = "http://localhost:8082"
	_ = productsURL
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{id}", displayCheckout)
	http.ListenAndServe(":8084", router)
}

func displayCheckout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Olá!")
}
