package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/bff/domain"

	outbound "github.com/nautible/nautible-app-ms-payment/pkg/generate/paymentclient"
)

type PaymentCashRepository struct{}

func NewPaymentCashRepository() domain.CashRepository {
	paymentCashRepository := PaymentCashRepository{}
	return &paymentCashRepository
}

// 決済登録を行うリポジトリ
func (p *PaymentCashRepository) CreatePayment(ctx context.Context, request *domain.Payment) (*domain.Payment, error) {
	c, err := outbound.NewClient("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-cash/method")
	if err != nil {
		panic(err)
	}
	var restCreatePayment outbound.RestCreatePayment
	restCreatePayment.CustomerId = request.CustomerId
	restCreatePayment.OrderDate = request.OrderDate
	restCreatePayment.OrderNo = request.OrderNo
	restCreatePayment.PaymentType = request.PaymentType
	restCreatePayment.RequestId = request.RequestId
	restCreatePayment.TotalPrice = request.TotalPrice
	requestJson, err := json.Marshal(restCreatePayment)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	buf := bytes.NewBuffer(requestJson)

	res, err := c.CreateWithBody(ctx, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 || res.StatusCode == 201 {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		var result domain.Payment
		json.Unmarshal(buf.Bytes(), &result)
		return &result, nil
	}
	return nil, errors.New("http response not success")
}

// 決済データ取得を行うリポジトリ
func (p *PaymentCashRepository) GetByOrderNo(ctx context.Context, orderNo string) (*domain.Payment, error) {
	fmt.Println("Rest GetByOrderNo")
	c, err := outbound.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-cash/method")
	if err != nil {
		return &domain.Payment{}, err
	}

	// http.Response として返却
	res, err := c.GetByOrderNoWithResponse(ctx, orderNo)
	if err != nil {
		return &domain.Payment{}, err
	}
	var model domain.Payment
	json.NewDecoder(bytes.NewReader(res.Body)).Decode(&model)
	return &model, err
}

// 決済データの取り消し
func (p *PaymentCashRepository) DeleteByOrderNo(ctx context.Context, orderNo string) error {
	fmt.Println("Rest DeleteByOrderNo")
	c, err := outbound.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-cash/method")
	if err != nil {
		return err
	}

	// http.Response として返却
	res, err := c.Delete(ctx, orderNo)
	if err != nil {
		return err
	}
	fmt.Println(res.StatusCode)
	return nil
}
