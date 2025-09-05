package dto

type CreateInvoiceRequest struct {
	TicketID string  `json:"ticket_id" binding:"required,uuid"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required,len=3"`
	DueDate  string  `json:"due_date" binding:"required"`
}


type CreatePaymentRequest struct {
	InvoiceID string  `json:"invoice_id" binding:"required,uuid"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Method    string  `json:"method" binding:"required,oneof=CREDIT_CARD BANK_TRANSFER PAYPAL"`
}

