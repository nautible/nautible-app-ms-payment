package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	domain "github.com/nautible/nautible-app-ms-payment/pkg/domain"
	"go.uber.org/zap"
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
		zap.S().Errorw("JSON Marshal error : " + err.Error())
		return err
	}
	buf := bytes.NewBuffer(requestJson)
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
