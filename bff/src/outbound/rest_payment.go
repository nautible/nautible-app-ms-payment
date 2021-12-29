package outbound

import (
	"bff/generate/client/payment"
	"bff/src/domain"
	"context"
	"net/http"
)

type paymentStruct struct{}

func NewRestPayment() domain.PaymentRepository {
	return &paymentStruct{}
}

func (p *paymentStruct) CreatePayment(ctx context.Context, request *http.Request) error {
	c, err := payment.NewClient("http://localhost:3500/v1.0/invoke/nautible-app-payment-cash/method/")
	if err != nil {
		panic(err)
	}

	// http.Response として返却
	res, err := c.CreateWithBody(context.Background(), "application/json", request.Body)
	defer res.Body.Close()
	return err
}

func (p *paymentStruct) GetByPaymentNo(ctx context.Context, paymentNo string) domain.PaymentItem {
	c, err := payment.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-payment-cash/method/")
	if err != nil {
		panic(err)
	}

	// http.Response として返却
	res, err := c.GetByPaymentNoWithResponse(context.Background(), paymentNo)
	if err != nil {
		panic(err)
	}

	return domain.PaymentItem{
		PaymentNo:   *res.JSON200.PaymentNo,
		AcceptNo:    *res.JSON200.AcceptNo,
		ReceiptDate: *res.JSON200.ReceiptDate,
		OrderNo:     *res.JSON200.OrderNo,
		OrderDate:   *res.JSON200.OrderDate,
		CustomerId:  *res.JSON200.CustomerId,
		TotalPrice:  *res.JSON200.TotalPrice,
		OrderStatus: *res.JSON200.OrderStatus,
	}
}
