// Package paymentserver provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package paymentserver

// RestCreatePayment defines model for RestCreatePayment.
type RestCreatePayment struct {
	CustomerId  int32  `json:"customerId"`
	OrderDate   string `json:"orderDate"`
	OrderNo     string `json:"orderNo"`
	PaymentType string `json:"paymentType"`

	// 処理要求を一意に表すリクエストId。
	RequestId  string `json:"requestId"`
	TotalPrice int32  `json:"totalPrice"`
}

// RestRejectCreatePayment defines model for RestRejectCreatePayment.
type RestRejectCreatePayment struct {
	// 処理要求を一意に表すリクエストId。
	OrderNo string `json:"orderNo"`
}

// CreateJSONBody defines parameters for Create.
type CreateJSONBody RestCreatePayment

// RejectCreateJSONBody defines parameters for RejectCreate.
type RejectCreateJSONBody RestRejectCreatePayment

// CreateJSONRequestBody defines body for Create for application/json ContentType.
type CreateJSONRequestBody CreateJSONBody

// RejectCreateJSONRequestBody defines body for RejectCreate for application/json ContentType.
type RejectCreateJSONRequestBody RejectCreateJSONBody
