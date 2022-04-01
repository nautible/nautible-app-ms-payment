package domain

// 後方サービスにリクエストするリポジトリインターフェース
type CreditRepository interface {
	CreatePayment(*PaymentItem) (*PaymentItem, error)
	GetByPaymentNo(string) (*PaymentItem, error)
}
