package inbound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	server "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditserver"
	dynamodb "github.com/nautible/nautible-app-ms-payment/pkg/outbound/dynamodb"
)

type CreditController struct {
	svc               *domain.CreditService
	RestPayment       server.RestCreditPayment
	RestUpdatePayment server.RestUpdateCreditPayment
	Lock              sync.Mutex
}

// Make sure we conform to ServerInterface

func NewCreditController(svc *domain.CreditService) *CreditController {
	return &CreditController{svc: svc}
}

// Healthz request
func (p *CreditController) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Health Check OK")
}

// Create Credit
// (POST /credit)
func (p *CreditController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req server.RestCreateCreditPayment
	json.NewDecoder(r.Body).Decode(&req)

	// サービス呼び出し
	var model domain.CreditPayment
	model.OrderNo = req.OrderNo
	model.OrderDate = req.OrderDate
	model.CustomerId = req.CustomerId
	model.TotalPrice = req.TotalPrice
	res, err := p.svc.CreateCreditPayment(r.Context(), &model)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var result server.RestCreditPayment
	result.AcceptNo = &res.AcceptNo
	result.AcceptDate = &res.AcceptDate
	result.CustomerId = &res.CustomerId
	result.OrderDate = &res.OrderDate
	result.OrderNo = &res.OrderNo
	result.TotalPrice = &res.TotalPrice
	resultJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

// Update Payment
// (PUT /payment/)
func (p *CreditController) Update(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(p.RestUpdatePayment)
	fmt.Fprint(w, string("Update No Operation"))
}

// Delete credit by AcceptNo
// (DELETE /credit/{acceptNo})
func (p *CreditController) Delete(w http.ResponseWriter, r *http.Request, acceptNo string) {

	repo := dynamodb.NewCreditRepository()
	svc := domain.NewCreditService(&repo)
	err := svc.DeleteCreditPayment(r.Context(), acceptNo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, string("Delete : "+acceptNo))
}

// Find credit by AcceptNo
// (GET /credit/{acceptNo})
func (p *CreditController) GetByAcceptNo(w http.ResponseWriter, r *http.Request, acceptNo string) {
	result, err := p.svc.GetCreditPayment(r.Context(), acceptNo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resultJson, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}
