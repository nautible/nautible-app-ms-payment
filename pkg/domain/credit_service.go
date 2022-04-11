package domain

import (
	"context"
	"fmt"
	"time"
)

type Credit struct {
	repo DbRepository
}

type CreditService interface {
	CreatePayment(context.Context, *PaymentModel) (*PaymentModel, error)
	GetPayment(context.Context, string) (*PaymentModel, error)
	DeletePayment(context.Context, string) error
}

func NewCreditService(repo DbRepository) CreditService {
	return &Credit{repo}
}

func (svc *Credit) CreatePayment(ctx context.Context, model *PaymentModel) (*PaymentModel, error) {
	paymentNo, err := svc.repo.Sequence(ctx)
	if err != nil {
		return nil, err
	}
	model.PaymentNo = fmt.Sprintf("C%010d", *paymentNo) // dummy クレジットの支払い番号はC始まりとする
	model.AcceptNo = fmt.Sprintf("A%010d", *paymentNo)  // dummy 受付番号はA始まりとする
	model.ReceiptDate = time.Now().String()             // dummy	return svc.repo.PutPaymentItem(ctx, model)
	return svc.repo.PutPaymentItem(ctx, model)
}

func (svc *Credit) GetPayment(ctx context.Context, orderNo string) (*PaymentModel, error) {
	return svc.repo.GetPaymentItem(ctx, orderNo)
}

func (svc *Credit) DeletePayment(ctx context.Context, orderNo string) error {
	return svc.repo.DeletePaymentItem(ctx, orderNo)
}
