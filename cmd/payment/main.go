package main

import (
	"log"
	"net/http"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	controller "github.com/nautible/nautible-app-ms-payment/pkg/inbound"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
	rest "github.com/nautible/nautible-app-ms-payment/pkg/outbound/rest"
)

func main() {
	controller := createController()
	http.HandleFunc("/payment/", func(w http.ResponseWriter, r *http.Request) {
		(*controller).IndexHandler(w, r)
	})
	http.HandleFunc("/payment/create", func(w http.ResponseWriter, r *http.Request) {
		(*controller).CreateHandler(w, r)
	})
	http.HandleFunc("/payment/rejectCreate", func(w http.ResponseWriter, r *http.Request) {
		(*controller).RejectCreateHandler(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createController() *controller.PaymentController {
	dynamoDbRepository := dynamodb.NewDynamoDbRepository()
	creditRepository := rest.NewCreditMessage()
	orderRepository := rest.NewOrderMessage()
	service := domain.NewPaymentService(&dynamoDbRepository, &creditRepository, &orderRepository)
	controller := controller.NewPaymentController(&service)
	return controller
}
