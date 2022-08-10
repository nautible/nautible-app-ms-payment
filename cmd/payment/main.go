package main

import (
	"log"
	"net/http"
	"os"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	controller "github.com/nautible/nautible-app-ms-payment/pkg/inbound"
	cosmosdb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/cosmosdb"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
	rest "github.com/nautible/nautible-app-ms-payment/pkg/outbound/rest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var target string // -ldflags '-X main.target=(aws|azure)'

func main() {
	var logger *zap.Logger
	var err error
	if os.Getenv("LOG_ENV") == "Development" {
		logger, err = NewDevelopmentLogger()
	} else {
		logger, err = NewProductionLogger()
	}
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	controller, repo := createController(target)
	defer (*repo).Close()

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

func createController(target string) (*controller.PaymentController, *domain.PaymentRepository) {
	var repo domain.PaymentRepository
	switch target {
	case "aws":
		repo = dynamodb.NewPaymentRepository()
	case "azure":
		repo = cosmosdb.NewPaymentRepository()
	default:
		panic("invalid ldflags parameter [main.target]")
	}
	creditMessage := rest.NewCreditMessageSender()
	orderMessage := rest.NewOrderMessageSender()
	service := domain.NewPaymentService(&repo, &creditMessage, &orderMessage)
	controller := controller.NewPaymentController(service)
	return controller, &repo
}

func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:   "level",
			TimeKey:    "timestamp",
			CallerKey:  "caller",
			MessageKey: "msg",
		},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func NewProductionLogger() (*zap.Logger, error) {
	config := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:   "level",
			TimeKey:    "timestamp",
			CallerKey:  "caller",
			MessageKey: "msg",
		},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
