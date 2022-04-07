package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/credit/domain"
	outbound "github.com/nautible/nautible-app-ms-payment/pkg/credit/outbound"
	server "github.com/nautible/nautible-app-ms-payment/pkg/generate/backendserver"
)

type Payment struct {
	RestPayment       server.RestPayment
	RestUpdatePayment server.RestUpdatePayment
	Lock              sync.Mutex
}

// Make sure we conform to ServerInterface

var _ server.ServerInterface = (*Payment)(nil)

func NewPayment() *Payment {
	return &Payment{}
}

// Find payments
// (GET /payment)
func (p *Payment) Find(w http.ResponseWriter, r *http.Request, params server.FindParams) {
	fmt.Fprint(w, string("Find"))
}

// Create Payment
// (POST /payment)
func (p *Payment) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	repo := outbound.NewPaymentDB()
	svc := domain.NewPaymentService(repo)
	var req server.RestCreatePayment
	json.NewDecoder(r.Body).Decode(&req)

	// サービス呼び出し
	var payment domain.Payment
	payment.CustomerId = req.CustomerId
	payment.TotalPrice = req.TotalPrice
	payment.OrderDate = req.OrderDate
	payment.OrderNo = req.OrderNo
	res, err := svc.CreatePayment(r.Context(), &payment)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var result server.RestPayment
	result.AcceptNo = &res.AcceptNo
	result.CustomerId = &res.CustomerId
	result.OrderDate = &res.OrderDate
	result.OrderNo = &res.OrderNo
	result.OrderStatus = &res.OrderStatus
	result.PaymentNo = &res.PaymentNo
	result.ReceiptDate = &res.ReceiptDate
	result.TotalPrice = &res.TotalPrice
	result.RequestId = &req.RequestId
	resultJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

// Update Payment
// (PUT /payment/)
func (p *Payment) Update(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.RestUpdatePayment)
	fmt.Fprint(w, string("Update"))
}

// Delete payment by orderNo
// (DELETE /payment/{orderNo})
func (p *Payment) Delete(w http.ResponseWriter, r *http.Request, orderNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")
	fmt.Fprint(w, string("Delete : "+id))

	repo := outbound.NewPaymentDB()
	svc := domain.NewPaymentService(repo)
	err := svc.DeletePayment(r.Context(), id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Find order by OrderNo
// (GET /payment/{orderNo})
func (p *Payment) GetByOrderNo(w http.ResponseWriter, r *http.Request, orderNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")

	repo := outbound.NewPaymentDB()
	svc := domain.NewPaymentService(repo)
	result, err := svc.GetPayment(r.Context(), id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}
