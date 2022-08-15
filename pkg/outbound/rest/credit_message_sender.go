package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.uber.org/zap"

	client "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditclient"
)

type CreditMessageSender struct{}

func NewCreditMessageSender() domain.CreditMessage {
	creditMessageSender := CreditMessageSender{}
	return &creditMessageSender
}

// クレジット決済登録を依頼するメッセージ
func (p *CreditMessageSender) CreateCreditPayment(ctx context.Context, request *domain.CreditPayment) (*domain.CreditPayment, error) {
	c, err := client.NewClient("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		panic(err)
	}
	var restCreatePayment client.RestCreateCreditPayment
	restCreatePayment.CustomerId = request.CustomerId
	restCreatePayment.OrderDate = request.OrderDate
	restCreatePayment.OrderNo = request.OrderNo
	restCreatePayment.TotalPrice = request.TotalPrice
	requestJson, err := json.Marshal(restCreatePayment)
	if err != nil {
		zap.S().Errorw("JSON Marshal error : " + err.Error())
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
		var result domain.CreditPayment
		json.Unmarshal(buf.Bytes(), &result)
		return &result, nil
	}
	return nil, errors.New("http response not success")
}

// クレジット決済データ取得を依頼するメッセージ
func (p *CreditMessageSender) GetByAcceptNo(ctx context.Context, acceptNo string) (*domain.CreditPayment, error) {
	c, err := client.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		return &domain.CreditPayment{}, err
	}

	// http.Response として返却
	res, err := c.GetByAcceptNoWithResponse(ctx, acceptNo)
	if err != nil {
		return &domain.CreditPayment{}, err
	}
	var model domain.CreditPayment
	json.NewDecoder(bytes.NewReader(res.Body)).Decode(&model)
	return &model, err
}

// 決済データの取り消し
func (p *CreditMessageSender) DeleteByAcceptNo(ctx context.Context, acceptNo string) error {
	c, err := client.NewClientWithResponses("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method")
	if err != nil {
		return err
	}

	// http.Response として返却
	res, err := c.Delete(ctx, acceptNo)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 204 {
		return errors.New("StatusCode is not 204 StatusCode:" + strconv.Itoa(res.StatusCode))
	}
	return nil
}
