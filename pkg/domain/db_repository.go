package domain

import "context"

// PaymentHistoryテーブルにアクセスするリポジトリインターフェース
type DbRepository interface {
	PutPaymentItem(context.Context, *PaymentModel) (*PaymentModel, error)
	GetPaymentItem(context.Context, string) (*PaymentModel, error)
	DeletePaymentItem(context.Context, string) error
	PutPaymentHistory(context.Context, *PaymentModel) error
	Sequence(context.Context) (*int, error)
}
