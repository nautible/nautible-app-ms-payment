package domain

import "context"

// PaymentHistoryテーブルにアクセスするリポジトリインターフェース
type DynamoDbRepository interface {
	PutPaymentHistory(context.Context, *Payment) error
}
