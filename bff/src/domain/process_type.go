package domain

type ProcessType string

const (
	StockReserveAllocate = ProcessType("STOCK_RESERVE_ALLOCATE")
	StockApproveAllocate = ProcessType("STOCK_APPROVE_ALLOCATE")
	StocRejectAllocate   = ProcessType("STOCK_REJECT_ALLOCATE")
	PaymentCreate        = ProcessType("PAYMENT_CREATE")
	PaymentRejectCreate  = ProcessType("PAYMENT_REJECT_CREATE")
	Payment              = ProcessType("PAYMENT")
)
