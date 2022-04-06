package domain

import "context"

// 後方サービスにリクエストするリポジトリインターフェース
type CashRepository interface {
	CreatePayment(context.Context, *Payment) (*Payment, error)
	GetByOrderNo(context.Context, string) (*Payment, error)
	DeleteByOrderNo(context.Context, string) error
}
