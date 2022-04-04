package domain

type PaymentType string

const (
	Credit = PaymentType("01")
	Cash   = PaymentType("02")
)
