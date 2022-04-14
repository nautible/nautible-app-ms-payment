package domain

import "context"

// 後方サービスにリクエストするメッセージインターフェース
type CreditMessageSender interface {
	Find(context.Context, int32, string, string) ([]*Payment, error)
	Create(context.Context, *Payment) (*Payment, error)
	GetByOrderNo(context.Context, string) (*Payment, error)
	DeleteByOrderNo(context.Context, string) error
}
