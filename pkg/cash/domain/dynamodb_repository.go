package domain

import (
	"context"
)

type DynamoDbRepository interface {
	PutPaymentItem(context.Context, *Payment) (*Payment, error)
	GetPaymentItem(context.Context, string) (*Payment, error)
	DeletePaymentItem(context.Context, string) error
}
