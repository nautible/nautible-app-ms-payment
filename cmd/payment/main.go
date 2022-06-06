package main

import (
	"fmt"
	"log"
	"net/http"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	controller "github.com/nautible/nautible-app-ms-payment/pkg/inbound"
	cosmosdb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/cosmosdb"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
	rest "github.com/nautible/nautible-app-ms-payment/pkg/outbound/rest"
)

var target string // -ldflags '-X main.target=(aws|azure)'

func main() {
	controller := createController(target)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		controller.HealthCheck(w, r)
	})
	http.HandleFunc("/payment/create", func(w http.ResponseWriter, r *http.Request) {
		controller.Create(w, r)
	})
	http.HandleFunc("/payment/rejectCreate", func(w http.ResponseWriter, r *http.Request) {
		controller.RejectCreate(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createController(target string) *controller.PaymentController {
	var repo domain.PaymentRepository
	switch target {
	case "aws":
		fmt.Println("dynamodb...")
		repo = dynamodb.NewPaymentRepository()
	case "azure":
		fmt.Println("cosmosdb...")
		repo = cosmosdb.NewPaymentRepository()
	default:
		panic("invalid ldflags parameter [main.target]")
	}
	creditMessage := rest.NewCreditMessageSender()
	orderMessage := rest.NewOrderMessageSender()
	service := domain.NewPaymentService(&repo, &creditMessage, &orderMessage)
	controller := controller.NewPaymentController(service)
	return controller
}
