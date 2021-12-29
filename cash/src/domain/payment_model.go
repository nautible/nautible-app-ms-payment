package domain

type PaymentItem struct {
	PaymentNo   string
	AcceptNo    string
	ReceiptDate string
	OrderNo     string
	OrderDate   string
	CustomerId  int32
	TotalPrice  int32
	OrderStatus string
}
