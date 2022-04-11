package inbound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	server "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditserver"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
)

type CreditController struct {
	svc               *domain.CreditService
	RestPayment       server.RestPayment
	RestUpdatePayment server.RestUpdatePayment
	Lock              sync.Mutex
}

// Make sure we conform to ServerInterface

var _ server.ServerInterface = (*CreditController)(nil)

func NewCreditController(svc *domain.CreditService) *CreditController {
	return &CreditController{svc: svc}
}

// Find payments
// (GET /payment)
func (p *CreditController) Find(w http.ResponseWriter, r *http.Request, params server.FindParams) {
	fmt.Fprint(w, string("Find"))
}

// Create Payment
// (POST /payment)
func (p *CreditController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req server.RestCreatePayment
	json.NewDecoder(r.Body).Decode(&req)

	// サービス呼び出し
	var model domain.PaymentModel
	model.CustomerId = req.CustomerId
	model.TotalPrice = req.TotalPrice
	model.OrderDate = req.OrderDate
	model.OrderNo = req.OrderNo
	service := *p.svc
	res, err := service.CreatePayment(r.Context(), &model)
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
func (p *CreditController) Update(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.RestUpdatePayment)
	fmt.Fprint(w, string("Update"))
}

// Delete payment by orderNo
// (DELETE /payment/{orderNo})
func (p *CreditController) Delete(w http.ResponseWriter, r *http.Request, orderNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")
	fmt.Fprint(w, string("Delete : "+id))

	repo := dynamodb.NewDynamoDbRepository()
	svc := domain.NewCreditService(repo)
	err := svc.DeletePayment(r.Context(), id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Find order by OrderNo
// (GET /payment/{orderNo})
func (p *CreditController) GetByOrderNo(w http.ResponseWriter, r *http.Request, orderNo string) {
	id := strings.TrimPrefix(r.URL.Path, "/payment/")

	repo := dynamodb.NewDynamoDbRepository()
	svc := domain.NewCreditService(repo)
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
