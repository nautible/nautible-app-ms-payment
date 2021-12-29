package domain

import (
	"context"
	"net/http"
)

// 後方サービスにリクエストするリポジトリインターフェース
type PaymentRepository interface {
	CreatePayment(context.Context, *http.Request) error
	GetByPaymentNo(context.Context, string) PaymentItem
}
