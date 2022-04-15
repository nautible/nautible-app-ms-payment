package domain

import "context"

// 後方サービスにリクエストするメッセージインターフェース
type CreditMessageService interface {
	CreateCreditPayment(context.Context, *CreditPayment) (*CreditPayment, error)
	GetByAcceptNo(context.Context, string) (*CreditPayment, error)
	DeleteByAcceptNo(context.Context, string) error
}
