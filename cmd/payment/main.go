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
	http.HandleFunc("/helthz", func(w http.ResponseWriter, r *http.Request) {
		controller.HelthCheck(w, r)
	})
	http.HandleFunc("/payment/create", func(w http.ResponseWriter, r *http.Request) {
		controller.Create(w, r)
	})
	http.HandleFunc("/payment/rejectCreate", func(w http.ResponseWriter, r *http.Request) {
		controller.RejectCreate(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createController() *controller.PaymentController {
	paymentRepository := dynamodb.NewPaymentRepository()
	creditMessage := rest.NewCreditMessageSender()
	orderMessage := rest.NewOrderMessageSender()
	service := domain.NewPaymentService(&paymentRepository, &creditMessage, &orderMessage)
	controller := controller.NewPaymentController(service)
	return controller
}
