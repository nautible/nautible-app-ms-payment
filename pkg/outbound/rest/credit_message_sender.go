package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"

	client "github.com/nautible/nautible-app-ms-payment/pkg/generate/creditclient"
)

type CreditMessageSender struct{}

func NewCreditMessageSender() domain.CreditMessage {
	creditMessageSender := CreditMessageSender{}
	return &creditMessageSender
}

// クレジット決済登録を依頼するメッセージ
func (p *CreditMessageSender) CreateCreditPayment(ctx context.Context, request *domain.CreditPayment) (*domain.CreditPayment, error) {
	header := ctx.Value("header").(http.Header)

	var restCreatePayment client.RestCreateCreditPayment
	restCreatePayment.CustomerId = request.CustomerId
	restCreatePayment.OrderDate = request.OrderDate
	restCreatePayment.OrderNo = request.OrderNo
	restCreatePayment.TotalPrice = request.TotalPrice
	requestJson, err := json.Marshal(restCreatePayment)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	buf := bytes.NewBuffer(requestJson)
	req, err := client.NewCreateRequestWithBody("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method", "application/json", buf)
	if err != nil {
		return nil, err
	}
	// トレーシング用ヘッダを伝搬させる
	req.Header.Add("x-b3-traceid", header.Get("x-b3-traceid"))
	req.Header.Add("x-b3-spanid", header.Get("x-b3-spanid"))
	req.Header.Add("x-b3-parentspanid", header.Get("x-b3-parentspanid"))
	req.Header.Add("x-b3-sampled", header.Get("x-b3-sampled"))
	req.Header.Add("x-b3-flags", header.Get("x-b3-flags"))
	client := &http.Client{}
	res, err := client.Do(req)
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
	header := ctx.Value("header").(http.Header)
	// http.Response として返却
	req, err := client.NewGetByAcceptNoRequest("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method", acceptNo)
	if err != nil {
		return nil, err
	}
	// トレーシング用ヘッダを伝搬させる
	req.Header.Add("x-b3-traceid", header.Get("x-b3-traceid"))
	req.Header.Add("x-b3-spanid", header.Get("x-b3-spanid"))
	req.Header.Add("x-b3-parentspanid", header.Get("x-b3-parentspanid"))
	req.Header.Add("x-b3-sampled", header.Get("x-b3-sampled"))
	req.Header.Add("x-b3-flags", header.Get("x-b3-flags"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &domain.CreditPayment{}, err
	}
	defer res.Body.Close()
	var model domain.CreditPayment
	json.NewDecoder(res.Body).Decode(&model)
	return &model, err
}

// 決済データの取り消し
func (p *CreditMessageSender) DeleteByAcceptNo(ctx context.Context, acceptNo string) error {
	header := ctx.Value("header").(http.Header)
	// http.Response として返却
	req, err := client.NewDeleteRequest("http://localhost:3500/v1.0/invoke/nautible-app-ms-payment-credit/method", acceptNo)
	if err != nil {
		return err
	}
	// トレーシング用ヘッダを伝搬させる
	req.Header.Add("x-b3-traceid", header.Get("x-b3-traceid"))
	req.Header.Add("x-b3-spanid", header.Get("x-b3-spanid"))
	req.Header.Add("x-b3-parentspanid", header.Get("x-b3-parentspanid"))
	req.Header.Add("x-b3-sampled", header.Get("x-b3-sampled"))
	req.Header.Add("x-b3-flags", header.Get("x-b3-flags"))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return errors.New("StatusCode is not 204 StatusCode:" + strconv.Itoa(res.StatusCode))
	}
	return nil
}
