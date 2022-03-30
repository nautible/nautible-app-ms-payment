package domain

import (
	"context"
)

type PaymentService interface {
	CreatePayment(context.Context, *PaymentItem) (*PaymentItem, error)
	GetPayment(context.Context, string) (*PaymentItem, error)
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (svc *paymentService) CreatePayment(ctx context.Context, payment *PaymentItem) (*PaymentItem, error) {
	return svc.repo.PutItem(ctx, payment)
}

func (svc *paymentService) GetPayment(ctx context.Context, paymentNo string) (*PaymentItem, error) {
	return svc.repo.GetItem(ctx, paymentNo)
}
