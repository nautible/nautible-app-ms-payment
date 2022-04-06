package domain

import "context"

// Orderサービスにリクエストするリポジトリインターフェース
type OrderRepository interface {
	PaymentResponse(context.Context, *OrderResponse) error
}
