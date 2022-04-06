package outbound

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	domain "payment-bff/src/domain"
)

type OrderRepository struct{}

func NewOrderRepository() domain.OrderRepository {
	orderRepository := OrderRepository{}
	return &orderRepository
}

// Orderサービスにリクエストするリポジトリインターフェース
func (p *OrderRepository) PaymentResponse(ctx context.Context, response *domain.OrderResponse) error {
	url := "http://localhost:3500/v1.0/publish/order-pubsub/create-order-reply"
	requestJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return err
	}
	buf := bytes.NewBuffer(requestJson)
	str := string(buf.Bytes())
	fmt.Println(str)
	// http.Response として返却
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	req.Header.Set("Content-Type", "application/octet-stream")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
