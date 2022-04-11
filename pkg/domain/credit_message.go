package domain

import "context"

// 後方サービスにリクエストするメッセージインターフェース
type CreditMessage interface {
	CreatePayment(context.Context, *PaymentModel) (*PaymentModel, error)
	GetByOrderNo(context.Context, string) (*PaymentModel, error)
	DeleteByOrderNo(context.Context, string) error
}
