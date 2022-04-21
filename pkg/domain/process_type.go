package domain

type ProcessType string

const (
	TypeStockReserveAllocate = ProcessType("STOCK_RESERVE_ALLOCATE")
	TypeStockApproveAllocate = ProcessType("STOCK_APPROVE_ALLOCATE")
	TypeStocRejectAllocate   = ProcessType("STOCK_REJECT_ALLOCATE")
	TypePaymentCreate        = ProcessType("PAYMENT_CREATE")
	TypePaymentRejectCreate  = ProcessType("PAYMENT_REJECT_CREATE")
	TypePayment              = ProcessType("PAYMENT")
)
