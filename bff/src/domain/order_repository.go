package domain

// Orderサービスにリクエストするリポジトリインターフェース
type OrderRepository interface {
	PaymentResponse(*OrderResponse) error
}
