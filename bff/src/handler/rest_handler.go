package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	server "payment-bff/generate/server"
	domain "payment-bff/src/domain"
	outbound "payment-bff/src/outbound"
)

type CloudEvents struct {
	Pubsubname      string `json:"pubsubname"`
	Id              string `json:"id"`
	Specversion     string `json:"specversion"`
	Source          string `json:"source"`
	Topic           string `json:"topic"`
	Datacontenttype string `json:"datacontenttype"`
	Type            string `json:"type"`
	Traceid         string `json:"traceid"`
	DataBase64      string `json:"data_base64"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, IndexHandler")
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateHandler")
	fmt.Println(r.Method)
	switch r.Method {
	case "POST":
		postCreate(w, r)
	default:
		log.Fatalf("%s Method not allowed.\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RejectCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postRejectCreate(w, r)
	default:
		log.Fatalf("%s Method not allowed.\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("postCreate")
	w.WriteHeader(http.StatusCreated)
	body := r.Body
	defer body.Close()

	// CloudEventsで受け取ったバイナリデータ（Base64）を構造体にマッピング
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	var cloudEvents CloudEvents
	var restCreatePayment server.RestCreatePayment
	json.Unmarshal(buf.Bytes(), &cloudEvents)
	dec, err := base64.StdEncoding.DecodeString(cloudEvents.DataBase64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.Unmarshal(dec, &restCreatePayment)

	// 決済サービス
	payment := outbound.NewPaymentRepository()
	order := outbound.NewOrderRepository()
	service := domain.NewPaymentService(&payment, &order)
	var paymentItem domain.PaymentItem
	paymentItem.RequestId = restCreatePayment.RequestId
	paymentItem.OrderNo = restCreatePayment.OrderNo
	paymentItem.PaymentType = restCreatePayment.PaymentType
	paymentItem.OrderDate = restCreatePayment.OrderDate
	paymentItem.TotalPrice = restCreatePayment.TotalPrice
	service.CreatePayment(&paymentItem)
	// result, err := service.GetByPaymentNo("C000000001")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result.OrderDate)
	// result2, err := service.GetByPaymentNo("C000000002")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result2.OrderDate)
	w.WriteHeader(http.StatusOK)
}

func postRejectCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Hello, RejectCreateHandler")
}
