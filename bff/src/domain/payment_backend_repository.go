package domain

// 後方サービスにリクエストするリポジトリインターフェース
type PaymentRepository interface {
	CreatePayment(*PaymentItem) (*PaymentItem, error)
	GetByPaymentNo(string) (*PaymentItem, error)
}
