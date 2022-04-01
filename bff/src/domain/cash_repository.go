package domain

// 後方サービスにリクエストするリポジトリインターフェース
type CashRepository interface {
	CreatePayment(*PaymentItem) (*PaymentItem, error)
	GetByPaymentNo(string) (*PaymentItem, error)
}
