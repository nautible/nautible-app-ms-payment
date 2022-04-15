package domain

// クレジット決済のモデル定義
type CreditPayment struct {
	AcceptNo   string `json:"acceptNo"`
	AcceptDate string `json:"acceptDate"`
	OrderNo    string `json:"orderNo"`
	OrderDate  string `json:"orderDate" validate:"required,datetime=2006-01-02T03:04:05"`
	CustomerId int32  `json:"customerId" validate:"required,gte=0"`
	TotalPrice int32  `json:"totalPrice" validate:"required,gte=0,lte=999999"`
	DeleteFlag bool   `json:"deleteFlag"`
}
