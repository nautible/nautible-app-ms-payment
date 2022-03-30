package domain

// Orderへのレスポンスのモデル定義
type OrderResponse struct {
	RequestId   string `json:"requestId"`
	Status      int    `json:"status"`
	ProcessType string `json:"processType"`
	Message     string `json:"message"`
}
