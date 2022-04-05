package domain

// 後方サービスにリクエストするリポジトリインターフェース
type CreditRepository interface {
	CreatePayment(*Payment) (*Payment, error)
	GetByOrderNo(string) (*Payment, error)
	DeleteByOrderNo(string) error
}
