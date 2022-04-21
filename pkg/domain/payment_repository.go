package domain

import "context"

// Paymentテーブルにアクセスするリポジトリインターフェース
type PaymentRepository interface {
	FindPayment(context.Context, int32, string, string) ([]*Payment, error)
	PutPayment(context.Context, *Payment) (*Payment, error)
	GetPayment(context.Context, string) (*Payment, error)
	DeletePayment(context.Context, string) error
	PutPaymentHistory(context.Context, *Payment) error
	Sequence(context.Context) (*int, error)
}
