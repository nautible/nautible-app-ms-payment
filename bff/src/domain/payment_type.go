package domain

type PaymentType string

const (
	Credit = ProcessType("01")
	Cash   = ProcessType("02")
)
