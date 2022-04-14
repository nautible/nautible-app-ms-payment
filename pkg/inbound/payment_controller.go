package inbound

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	server "github.com/nautible/nautible-app-ms-payment/pkg/generate/paymentserver"
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

type PaymentController struct {
	svc  *domain.PaymentService
	Lock sync.Mutex
}

var _ *PaymentController = (*PaymentController)(nil)

func NewPaymentController(svc *domain.PaymentService) *PaymentController {
	return &PaymentController{svc: svc}
}

func (p *PaymentController) HelthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Helth Check OK")
}

func (p *PaymentController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateHandler")
	fmt.Println(r.Method)
	switch r.Method {
	case "POST":
		doCreate(w, r, p.svc)
	default:
		log.Fatalf("%s Method not allowed.\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *PaymentController) RejectCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		doRejectCreate(w, r, p.svc)
	default:
		log.Fatalf("%s Method not allowed.\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// /create受信ハンドラ
func doCreate(w http.ResponseWriter, r *http.Request, svc *domain.PaymentService) {
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

	// 入力データの取得
	var model domain.PaymentModel
	model.RequestId = restCreatePayment.RequestId
	model.OrderNo = restCreatePayment.OrderNo
	model.PaymentType = restCreatePayment.PaymentType
	model.OrderDate = restCreatePayment.OrderDate
	model.TotalPrice = restCreatePayment.TotalPrice
	model.CustomerId = restCreatePayment.CustomerId

	// 決済サービス呼び出し
	(*svc).CreatePayment(r.Context(), &model)

	w.WriteHeader(http.StatusOK)
}

func doRejectCreate(w http.ResponseWriter, r *http.Request, svc *domain.PaymentService) {
	w.WriteHeader(http.StatusCreated)
	body := r.Body
	defer body.Close()

	// CloudEventsで受け取ったバイナリデータ（Base64）を構造体にマッピング
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	var cloudEvents CloudEvents
	var restCancelPayment server.RestCancelPayment
	json.Unmarshal(buf.Bytes(), &cloudEvents)
	dec, err := base64.StdEncoding.DecodeString(cloudEvents.DataBase64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.Unmarshal(dec, &restCancelPayment)

	// 決済削除サービス呼び出し
	(*svc).DeleteByOrderNo(r.Context(), restCancelPayment.PaymentType, restCancelPayment.OrderNo)

	w.WriteHeader(http.StatusOK)
}
