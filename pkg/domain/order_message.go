package domain

import "context"

// 他サービスに非同期リクエストするメッセージインターフェース
type OrderMessage interface {
	Publish(context.Context, interface{}) error
}
