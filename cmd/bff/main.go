package main

import (
	"log"
	"net/http"

	handler "github.com/nautible/nautible-app-ms-payment/pkg/bff/handler"
)

func inject() {

}
func main() {
	http.HandleFunc("/payment/", handler.IndexHandler)
	http.HandleFunc("/payment/create", handler.CreateHandler)
	http.HandleFunc("/payment/rejectCreate", handler.RejectCreateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
