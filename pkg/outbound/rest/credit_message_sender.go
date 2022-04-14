package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"

	client "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditclient"
)

type CreditMessage struct{}

func NewCreditMessage() domain.CreditMessageSender {
	creditMessage := CreditMessage{}
	return &creditMessage
}

func (p *CreditMessage) Find(ctx context.Context, customerId int32, orderDateFrom string, orderDateTo string) ([]*domain.Payment, error) {
	c, err := client.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		panic(err)
	}
	var param *client.FindParams
	param.CustomerId = &customerId
	param.OrderDateFrom = &orderDateFrom
	param.OrderDateTo = &orderDateTo

	res, err := c.FindWithResponse(ctx, param)
	if err != nil {
		return nil, err
	}

	if res.HTTPResponse.StatusCode == http.StatusOK || res.HTTPResponse.StatusCode == http.StatusCreated {
		var result []*domain.Payment
		err := json.Unmarshal(res.Body, &result)
		return result, err
	}
	return nil, errors.New("http response not success")
}

// 決済登録を行うリポジトリ
func (p *CreditMessage) Create(ctx context.Context, request *domain.Payment) (*domain.Payment, error) {
	c, err := client.NewClient("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		panic(err)
	}
	var restCreatePayment client.RestCreatePayment
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
func (p *CreditMessage) GetByOrderNo(ctx context.Context, paymentNo string) (*domain.Payment, error) {
	fmt.Println("Rest GetByPaymentNo")
	c, err := client.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
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
func (p *CreditMessage) DeleteByOrderNo(ctx context.Context, orderNo string) error {
	fmt.Println("Rest DeleteByOrderNo")
	c, err := client.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-cash/method")
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
