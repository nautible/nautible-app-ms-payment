package inbound

import (
	"bff/generate/server"
	"bff/src/domain"
	"bff/src/outbound"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Payment struct {
	Subscribe         server.DaprSubscribe
	RestPayment       server.RestPayment
	RestUpdatePayment server.RestUpdatePayment
	Lock              sync.Mutex
}

// Make sure we conform to ServerInterface

var _ server.ServerInterface = (*Payment)(nil)

func NewPayment() *Payment {
	return &Payment{}
}

// Here, we implement all of the handlers in the ServerInterface
// daprSubscribe
// (GET /dapr/subscribe)
func (p *Payment) DaprSubscribe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p.Subscribe)
	fmt.Fprint(w, string("DaprSubscribe"))
}

// Find payments
// (GET /payment/)
func (p *Payment) Find(w http.ResponseWriter, r *http.Request, params server.FindParams) {
	fmt.Fprint(w, string("Find"))

}

// Update Payment
// (PUT /payment/)
func (p *Payment) Update(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.RestUpdatePayment)
	fmt.Fprint(w, string("Update"))
}

// Create Payment for SAGA
// (POST /payment/create)
func (p *Payment) Create(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/create/")
	repo := outbound.NewRestPayment()

	svc := domain.NewPaymentService(repo)

	svc.CreatePayment(context.Background(), r)

	fmt.Fprint(w, string("Create : "+id))
}

// Reject Create Payment for SAGA
// (POST /payment/rejectCreate)
func (p *Payment) RejectCreate(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/rejectCreate/")
	fmt.Fprint(w, string("Create : "+id))
}

// Delete payment by paymentNo
// (DELETE /payment/{paymentNo})
func (p *Payment) Delete(w http.ResponseWriter, r *http.Request, paymentNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")
	fmt.Fprint(w, string("Delete : "+id))
}

// Find order by PaymentNo
// (GET /payment/{paymentNo})
func (p *Payment) GetByPaymentNo(w http.ResponseWriter, r *http.Request, paymentNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")
	fmt.Println("id : " + id)

	repo := outbound.NewRestPayment()

	svc := domain.NewPaymentService(repo)
	result := svc.GetByPaymentNo(context.Background(), id)
	fmt.Fprint(w, string("OrderNo : "+result.OrderNo+" AcceptNo : "+result.AcceptNo))
}
