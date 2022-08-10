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
	"go.uber.org/zap"
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
	Data            string `json:"data"`
	DataBase64      string `json:"data_base64"`
}

type PaymentController struct {
	svc  *domain.PaymentService
	Lock sync.Mutex
}

func NewPaymentController(svc *domain.PaymentService) *PaymentController {
	return &PaymentController{svc: svc}
}

func (p *PaymentController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Health Check OK")
}

func (p *PaymentController) Create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		doCreate(w, r, p.svc)
	default:
		log.Fatalf("%s Method not allowed.\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *PaymentController) RejectCreate(w http.ResponseWriter, r *http.Request) {
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

	// CloudEventsで受け取ったデータを構造体にマッピング
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	zap.S().Infow("request : " + buf.String())
	var cloudEvents CloudEvents
	var restCreatePayment server.RestCreatePayment
	json.Unmarshal(buf.Bytes(), &cloudEvents)
	if cloudEvents.Data != "" {
		dec := []byte(cloudEvents.Data)
		json.Unmarshal(dec, &restCreatePayment)
	}
	if cloudEvents.DataBase64 != "" {
		dec, err := base64.StdEncoding.DecodeString(cloudEvents.DataBase64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.Unmarshal(dec, &restCreatePayment)
	}

	// 入力データの取得
	var model domain.Payment
	model.RequestId = restCreatePayment.RequestId
	model.OrderNo = restCreatePayment.OrderNo
	model.PaymentType = restCreatePayment.PaymentType
	model.OrderDate = restCreatePayment.OrderDate
	model.TotalPrice = restCreatePayment.TotalPrice
	model.CustomerId = restCreatePayment.CustomerId

	// 決済サービス呼び出し
	svc.CreatePayment(r.Context(), &model)
	w.WriteHeader(http.StatusOK)
}

func doRejectCreate(w http.ResponseWriter, r *http.Request, svc *domain.PaymentService) {
	body := r.Body
	defer body.Close()

	// CloudEventsで受け取ったバイナリデータ（Base64）を構造体にマッピング
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	zap.S().Infow("request : " + buf.String())
	var cloudEvents CloudEvents
	var restRejectCreatePayment server.RestRejectCreatePayment
	json.Unmarshal(buf.Bytes(), &cloudEvents)
	if cloudEvents.Data != "" {
		dec := []byte(cloudEvents.Data)
		json.Unmarshal(dec, &restRejectCreatePayment)
	}
	if cloudEvents.DataBase64 != "" {
		dec, err := base64.StdEncoding.DecodeString(cloudEvents.DataBase64)
		if err != nil {
			zap.S().Errorw(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.Unmarshal(dec, &restRejectCreatePayment)
	}

	// 決済削除サービス呼び出し
	svc.DeleteByOrderNo(r.Context(), restRejectCreatePayment.OrderNo)
	w.WriteHeader(http.StatusOK)
}
