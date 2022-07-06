package domain

// Paymentのモデル定義
type Payment struct {
	RequestId   string `json:"requestId" validate:"required,len=11"`
	PaymentNo   string `json:"paymentNo"`
	PaymentType string `json:"paymentType" validate:"required"`
	OrderNo     string `json:"orderNo" validate:"required,len=11"`
	OrderDate   string `json:"orderDate" validate:"required,datetime=2006-01-02T15:04:05"`
	AcceptNo    string `json:"acceptNo"`
	AcceptDate  string `json:"acceptDate"`
	CustomerId  int32  `json:"customerId" validate:"required,gte=0"`
	TotalPrice  int32  `json:"totalPrice" validate:"required,gte=0,lte=999999"`
	OrderStatus string `json:"orderStatus"`
	DeleteFlag  bool   `json:"deleteFlag"`
}
