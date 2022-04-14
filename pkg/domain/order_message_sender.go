package domain

import "context"

// Orderサービスにリクエストするメッセージインターフェース
type OrderMessageSender interface {
	Publish(context.Context, *Order) error
}
