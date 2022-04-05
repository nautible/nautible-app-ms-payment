package domain

// 後方サービスにリクエストするリポジトリインターフェース
type CashRepository interface {
	CreatePayment(*Payment) (*Payment, error)
	GetByOrderNo(string) (*Payment, error)
	DeleteByOrderNo(string) error
}
