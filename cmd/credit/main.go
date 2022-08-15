package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	server "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditserver"
	controller "github.com/nautible/nautible-app-ms-payment/pkg/inbound"
	cosmosdb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/cosmosdb"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
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

	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := server.GetSwagger()
	if err != nil {
		zap.S().Errorw("Error loading swagger spec : " + err.Error())
		os.Exit(1)
	}

	swagger.Servers = nil

	paymentController, repo := createController(target)
	defer (*repo).Close()

	r := chi.NewRouter()

	r.Use(middleware.OapiRequestValidator(swagger))

	server.HandlerFromMux(paymentController, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}
	log.Fatal(s.ListenAndServe())
}

func createController(target string) (*controller.CreditController, *domain.CreditRepository) {
	var repo domain.CreditRepository
	switch target {
	case "aws":
		repo = dynamodb.NewCreditRepository()
	case "azure":
		repo = cosmosdb.NewCreditRepository()
	default:
		panic("invalid ldflags parameter [main.target]")
	}
	svc := domain.NewCreditService(&repo)

	creditController := controller.NewCreditController(svc)
	return creditController, &repo
}

func NewDevelopmentLogger() (*zap.Logger, error) {
	config := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:       "level",
			TimeKey:        "timestamp",
			CallerKey:      "caller",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
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
			LevelKey:       "level",
			TimeKey:        "timestamp",
			CallerKey:      "caller",
			MessageKey:     "msg",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
