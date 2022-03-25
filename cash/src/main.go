package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"payment-cash/generate/server"
	"payment-cash/src/inbound"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("start main..........")

	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := server.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	swagger.Servers = nil

	new_payment := inbound.NewPayment()

	r := chi.NewRouter()

	r.Use(middleware.OapiRequestValidator(swagger))

	server.HandlerFromMux(new_payment, r)

	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}
	fmt.Println("start server..........")
	log.Fatal(s.ListenAndServe())
}
