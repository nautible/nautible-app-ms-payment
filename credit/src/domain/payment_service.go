package domain

import (
	"context"
)

type PaymentService interface {
	CreatePayment(context.Context, *Payment) (*Payment, error)
	GetPayment(context.Context, string) (*Payment, error)
	DeletePayment(context.Context, string) error
}

type paymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (svc *paymentService) CreatePayment(ctx context.Context, payment *Payment) (*Payment, error) {
	return svc.repo.PutItem(ctx, payment)
}

func (svc *paymentService) GetPayment(ctx context.Context, orderNo string) (*Payment, error) {
	return svc.repo.GetItem(ctx, orderNo)
}

func (svc *paymentService) DeletePayment(ctx context.Context, orderNo string) error {
	return svc.repo.DeleteItem(ctx, orderNo)
}
