package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
)

type OrderMessageSender struct{}

func NewOrderMessageSender() domain.OrderMessage {
	orderMessageSender := OrderMessageSender{}
	return &orderMessageSender
}

// Orderサービスにリクエストするリポジトリインターフェース
func (p *OrderMessageSender) Publish(ctx context.Context, response interface{}) error {
	url := "http://localhost:3500/v1.0/publish/order-pubsub/create-order-reply"
	requestJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return err
	}
	buf := bytes.NewBuffer(requestJson)
	str := string(buf.String())
	fmt.Println(str)
	// http.Response として返却
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
