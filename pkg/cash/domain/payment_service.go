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
	repo DynamoDbRepository
}

func NewPaymentService(repo DynamoDbRepository) PaymentService {
	return &paymentService{repo}
}

func (svc *paymentService) CreatePayment(ctx context.Context, payment *Payment) (*Payment, error) {
	return svc.repo.PutPaymentItem(ctx, payment)
}

func (svc *paymentService) GetPayment(ctx context.Context, paymentNo string) (*Payment, error) {
	return svc.repo.GetPaymentItem(ctx, paymentNo)
}

func (svc *paymentService) DeletePayment(ctx context.Context, orderNo string) error {
	return svc.repo.DeletePaymentItem(ctx, orderNo)
}
