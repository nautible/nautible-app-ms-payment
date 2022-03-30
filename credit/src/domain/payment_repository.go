package domain

import (
	"context"
)

type PaymentRepository interface {
	PutItem(context.Context, *PaymentItem) (*PaymentItem, error)
	GetItem(context.Context, string) (*PaymentItem, error)
}
