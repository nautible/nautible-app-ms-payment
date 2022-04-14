package domain

import (
	"context"
	"fmt"
	"time"
)

type CreditService struct {
	payment PaymentRepository
}

func NewCreditService(repo PaymentRepository) *CreditService {
	return &CreditService{repo}
}

func (svc *CreditService) CreatePayment(ctx context.Context, model *Payment) (*Payment, error) {
	paymentNo, err := svc.payment.Sequence(ctx)
	if err != nil {
		return nil, err
	}
	model.PaymentNo = fmt.Sprintf("C%010d", *paymentNo) // dummy クレジットの支払い番号はC始まりとする
	model.AcceptNo = fmt.Sprintf("A%010d", *paymentNo)  // dummy 受付番号はA始まりとする
	model.ReceiptDate = time.Now().String()             // dummy	return svc.repo.PutPaymentItem(ctx, model)
	return svc.payment.PutPaymentItem(ctx, model)
}

func (svc *CreditService) Find(ctx context.Context, customerId int32, orderDateFrom string, orderDateTo string) ([]*Payment, error) {
	return svc.payment.FindPaymentItem(ctx, customerId, orderDateFrom, orderDateTo)
}

func (svc *CreditService) GetPayment(ctx context.Context, orderNo string) (*Payment, error) {
	return svc.payment.GetPaymentItem(ctx, orderNo)
}

func (svc *CreditService) DeletePayment(ctx context.Context, orderNo string) error {
	return svc.payment.DeletePaymentItem(ctx, orderNo)
}
