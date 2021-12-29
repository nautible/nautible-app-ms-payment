// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package server

// DaprSubscribe defines model for DaprSubscribe.
type DaprSubscribe struct {
	// pubsub
	Pubsubname *string `json:"pubsubname,omitempty"`

	// route
	Route *string `json:"route,omitempty"`

	// topic
	Topic *string `json:"topic,omitempty"`
}

// RestPayment defines model for RestPayment.
type RestPayment struct {
	AcceptNo    *string `json:"acceptNo,omitempty"`
	CustomerId  *int32  `json:"customerId,omitempty"`
	OrderDate   *string `json:"orderDate,omitempty"`
	OrderNo     *string `json:"orderNo,omitempty"`
	OrderStatus *string `json:"orderStatus,omitempty"`
	PaymentNo   *string `json:"paymentNo,omitempty"`
	ReceiptDate *string `json:"receiptDate,omitempty"`
	TotalPrice  *int32  `json:"totalPrice,omitempty"`
}

// RestUpdatePayment defines model for RestUpdatePayment.
type RestUpdatePayment struct {
	AcceptNo    *string `json:"acceptNo,omitempty"`
	CustomerId  int32   `json:"customerId"`
	OrderDate   string  `json:"orderDate"`
	OrderNo     string  `json:"orderNo"`
	PaymentNo   string  `json:"paymentNo"`
	PaymentType string  `json:"paymentType"`
	TotalPrice  int32   `json:"totalPrice"`
}

// FindParams defines parameters for Find.
type FindParams struct {
	// customerId
	CustomerId *int32 `json:"customerId,omitempty"`

	// order date from
	OrderDateFrom *string `json:"orderDateFrom,omitempty"`

	// order date to
	OrderDateTo *string `json:"orderDateTo,omitempty"`
}

// UpdateJSONBody defines parameters for Update.
type UpdateJSONBody RestUpdatePayment

// UpdateJSONRequestBody defines body for Update for application/json ContentType.
type UpdateJSONRequestBody UpdateJSONBody
