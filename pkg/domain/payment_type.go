package domain

type PaymentType string

const (
	TypeCredit = PaymentType("01")
	TypeCash   = PaymentType("02")
)
