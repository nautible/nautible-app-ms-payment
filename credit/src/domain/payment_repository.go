package domain

import (
	"context"
)

type PaymentRepository interface {
	PutItem(context.Context, *Payment) (*Payment, error)
	GetItem(context.Context, string) (*Payment, error)
	DeleteItem(context.Context, string) error
}
