package domain

// Paymentのモデル定義
type PaymentItem struct {
	RequestId   string `json:"requestId"`
	PaymentNo   string `json:"paymentNo"`
	PaymentType string `json:"paymentType"`
	AcceptNo    string `json:"acceptNo"`
	ReceiptDate string `json:"receiptDate"`
	OrderNo     string `json:"orderNo"`
	OrderDate   string `json:"orderDate"`
	CustomerId  int32  `json:"customerId"`
	TotalPrice  int32  `json:"totalPrice"`
	OrderStatus string `json:"orderStatus"`
}
