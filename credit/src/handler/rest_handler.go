package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	server "payment-credit/generate/server"
	domain "payment-credit/src/domain"
	outbound "payment-credit/src/outbound"
	"strings"
	"sync"
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
	var model domain.Payment
	json.NewDecoder(r.Body).Decode(&model)
	res, err := svc.CreatePayment(r.Context(), &model)
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
	result.RequestId = &model.RequestId
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
