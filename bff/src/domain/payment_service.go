package domain

import (
	"context"
	"net/http"
)

type PaymentService interface {
	CreatePayment(context.Context, *http.Request)
	GetByPaymentNo(context.Context, string) PaymentItem
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (svc *paymentService) CreatePayment(ctx context.Context, req *http.Request) {
	svc.repo.CreatePayment(ctx, req)
}

func (svc *paymentService) GetByPaymentNo(ctx context.Context, paymentNo string) PaymentItem {
	return svc.repo.GetByPaymentNo(ctx, paymentNo)
}
