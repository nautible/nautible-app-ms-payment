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
func (p *PaymentCreditRepository) CreatePayment(request *domain.PaymentItem) (*domain.PaymentItem, error) {
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

	res, err := c.CreateWithBody(context.Background(), "application/json", buf)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		buf := new(bytes.Buffer)
		io.Copy(buf, res.Body)
		var result domain.PaymentItem
		json.Unmarshal(buf.Bytes(), &result)
		return &result, nil
	}
	return nil, errors.New("http response not success")
}

// 決済データ取得を行うリポジトリ
func (p *PaymentCreditRepository) GetByPaymentNo(paymentNo string) (*domain.PaymentItem, error) {
	fmt.Println("Rest GetByPaymentNo")
	c, err := outbound.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		return &domain.PaymentItem{}, err
	}

	// http.Response として返却
	res, err := c.GetByPaymentNoWithResponse(context.Background(), paymentNo)
	if err != nil {
		return &domain.PaymentItem{}, err
	}
	var model domain.PaymentItem
	json.NewDecoder(bytes.NewReader(res.Body)).Decode(&model)
	return &model, err
}
