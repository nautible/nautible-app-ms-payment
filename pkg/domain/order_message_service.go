package domain

import "context"

// 他サービスに非同期リクエストするメッセージインターフェース
type OrderMessageService interface {
	Send(context.Context, interface{}) error
}
