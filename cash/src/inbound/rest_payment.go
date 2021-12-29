package inbound

import (
	"cash/generate/server/payment"
	"cash/src/domain"
	"cash/src/outbound"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Payment struct {
	Subscribe         payment.DaprSubscribe
	RestPayment       payment.RestPayment
	RestUpdatePayment payment.RestUpdatePayment
	Lock              sync.Mutex
}

// Make sure we conform to ServerInterface

var _ payment.ServerInterface = (*Payment)(nil)

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
func (p *Payment) Find(w http.ResponseWriter, r *http.Request, params payment.FindParams) {
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
	paymentItem := domain.PaymentItem{
		PaymentNo:   *p.RestPayment.PaymentNo,
		AcceptNo:    *p.RestPayment.AcceptNo,
		ReceiptDate: *p.RestPayment.ReceiptDate,
		OrderNo:     *p.RestPayment.OrderNo,
		OrderDate:   *p.RestPayment.OrderDate,
		CustomerId:  *p.RestPayment.CustomerId,
		TotalPrice:  *p.RestPayment.TotalPrice,
		OrderStatus: *p.RestPayment.OrderStatus,
	}

	// セッション情報のべた書きは暫定
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://payment-localstack:4566"),
		Credentials: credentials.NewStaticCredentials("test-key", "test-secret", ""),
	})
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(session)
	repo := outbound.NewPaymentDB(db)
	svc := domain.NewPaymentService(repo)
	svc.CreatePayment(context.Background(), paymentItem)

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

	// セッション情報のべた書きは暫定
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Endpoint:    aws.String("http://payment-localstack:4566"),
		Credentials: credentials.NewStaticCredentials("test-key", "test-secret", ""),
	})
	if err != nil {
		panic(err)
	}
	db := dynamodb.New(session)
	repo := outbound.NewPaymentDB(db)
	svc := domain.NewPaymentService(repo)
	result := svc.GetPayment(context.Background(), id)
	fmt.Fprint(w, string("OrderNo : "+result.OrderNo+" AcceptNo : "+result.AcceptNo))
}
