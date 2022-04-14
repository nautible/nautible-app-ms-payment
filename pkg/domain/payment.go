package domain

// Paymentのモデル定義
type Payment struct {
	RequestId   string `json:"requestId" validate:"required,len=11"`
	PaymentNo   string `json:"paymentNo"`
	PaymentType string `json:"paymentType" validate:"required"`
	AcceptNo    string `json:"acceptNo"`
	ReceiptDate string `json:"receiptDate"`
	OrderNo     string `json:"orderNo" validate:"required,len=11"`
	OrderDate   string `json:"orderDate" validate:"required,datetime=2006-01-02T03:04:05"`
	CustomerId  int32  `json:"customerId" validate:"required,gte=0"`
	TotalPrice  int32  `json:"totalPrice" validate:"required,gte=0,lte=999999"`
	OrderStatus string `json:"orderStatus"`
	DeleteFlag  bool   `json:"deleteFlag"`
}
