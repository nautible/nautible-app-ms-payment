package domain

import (
	"fmt"
)

type PaymentService interface {
	CreatePayment(*PaymentItem)
	GetByPaymentNo(string) (*PaymentItem, error)
}

type PaymentStruct struct {
	payment *PaymentRepository
	order   *OrderRepository
}

func NewPaymentService(payment *PaymentRepository, order *OrderRepository) PaymentService {
	return &PaymentStruct{payment, order}
}

// バックエンドサービスに支払作成処理を投げ、結果をOrderに返す
func (svc *PaymentStruct) CreatePayment(payementItem *PaymentItem) {
	fmt.Println("CreatePaymentService")
	var orderResponse OrderResponse
	res, err := (*svc.payment).CreatePayment(payementItem)
	// エラー発生
	if err != nil {
		orderResponse.ProcessType = string(Payment)
		orderResponse.RequestId = res.RequestId
		orderResponse.Status = 503
		(*svc.order).PaymentResponse(&orderResponse)
	}
	// 正常応答
	orderResponse.ProcessType = string(Payment)
	orderResponse.RequestId = res.RequestId
	orderResponse.Status = 201
	(*svc.order).PaymentResponse(&orderResponse)
}

func (svc *PaymentStruct) GetByPaymentNo(paymentNo string) (*PaymentItem, error) {
	return (*svc.payment).GetByPaymentNo(paymentNo)
}
