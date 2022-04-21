package domain

import (
	"context"
	"fmt"
	"time"
)

type CreditService struct {
	credit CreditRepository
}

func NewCreditService(repo *CreditRepository) *CreditService {
	return &CreditService{*repo}
}

func (svc *CreditService) CreateCreditPayment(ctx context.Context, model *CreditPayment) (*CreditPayment, error) {
	paymentNo, err := svc.credit.Sequence(ctx)
	if err != nil {
		return nil, err
	}
	model.AcceptNo = fmt.Sprintf("A%010d", *paymentNo) // dummy 受付番号はA始まりとする
	model.AcceptDate = time.Now().String()             // dummy
	return svc.credit.PutCreditPayment(ctx, model)
}

func (svc *CreditService) GetCreditPayment(ctx context.Context, acceptNo string) (*CreditPayment, error) {
	return svc.credit.GetCreditPayment(ctx, acceptNo)
}

func (svc *CreditService) DeleteCreditPayment(ctx context.Context, acceptNo string) error {
	return svc.credit.DeleteCreditPayment(ctx, acceptNo)
}
