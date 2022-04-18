package domain

import "context"

// 後方サービスにリクエストするメッセージインターフェース
type CreditMessage interface {
	CreateCreditPayment(context.Context, *CreditPayment) (*CreditPayment, error)
	GetByAcceptNo(context.Context, string) (*CreditPayment, error)
	DeleteByAcceptNo(context.Context, string) error
}
