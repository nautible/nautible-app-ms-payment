package domain

import "context"

// Creditテーブルにアクセスするリポジトリインターフェース
type CreditRepository interface {
	PutCreditPayment(context.Context, *CreditPayment) (*CreditPayment, error)
	GetCreditPayment(context.Context, string) (*CreditPayment, error)
	DeleteCreditPayment(context.Context, string) error
	Sequence(context.Context) (*int, error)
}
