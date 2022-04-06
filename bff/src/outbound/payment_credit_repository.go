package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	outbound "payment-bff/generate/client/payment"
	domain "payment-bff/src/domain"
)

type PaymentCreditRepository struct{}

func NewPaymentCreditRepository() domain.CreditRepository {
	paymentCreditRepository := PaymentCreditRepository{}
	return &paymentCreditRepository
}

// 決済登録を行うリポジトリ
func (p *PaymentCreditRepository) CreatePayment(ctx context.Context, request *domain.Payment) (*domain.Payment, error) {
	c, err := outbound.NewClient("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
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
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		var result domain.Payment
		json.Unmarshal(buf.Bytes(), &result)
		return &result, nil
	}
	return nil, errors.New("http response not success")
}

// 決済データ取得を行うリポジトリ
func (p *PaymentCreditRepository) GetByOrderNo(ctx context.Context, paymentNo string) (*domain.Payment, error) {
	fmt.Println("Rest GetByPaymentNo")
	c, err := outbound.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		return &domain.Payment{}, err
	}

	// http.Response として返却
	res, err := c.GetByOrderNoWithResponse(ctx, paymentNo)
	if err != nil {
		return &domain.Payment{}, err
	}
	var model domain.Payment
	json.NewDecoder(bytes.NewReader(res.Body)).Decode(&model)
	return &model, err
}

// 決済データの取り消し
func (p *PaymentCreditRepository) DeleteByOrderNo(ctx context.Context, orderNo string) error {
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
