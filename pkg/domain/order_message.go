package domain

import "context"

// Orderサービスにリクエストするメッセージインターフェース
type OrderMessage interface {
	Publish(context.Context, *OrderResponse) error
}
