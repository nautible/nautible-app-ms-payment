package domain

import "context"

// Payment,PaymentHistoryテーブルにアクセスするリポジトリインターフェース
type PaymentRepository interface {
	FindPaymentItem(context.Context, int32, string, string) ([]*Payment, error)
	PutPaymentItem(context.Context, *Payment) (*Payment, error)
	GetPaymentItem(context.Context, string) (*Payment, error)
	DeletePaymentItem(context.Context, string) error
	PutPaymentHistory(context.Context, *Payment) error
	Sequence(context.Context) (*int, error)
}
