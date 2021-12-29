package domain

import (
	"context"
)

type PaymentService interface {
	CreatePayment(context.Context, PaymentItem)
	GetPayment(context.Context, string) PaymentItem
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (svc *paymentService) CreatePayment(ctx context.Context, payment PaymentItem) {
	svc.repo.PutItem(ctx, payment)
}

func (svc *paymentService) GetPayment(ctx context.Context, paymentNo string) PaymentItem {
	return svc.repo.GetItem(ctx, paymentNo)
}
