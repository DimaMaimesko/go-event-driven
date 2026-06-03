package entities

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`

	IdempotencyKey string `json:"idempotency_key"`
}

type VoidReceipt struct {
	TicketID       string `json:"ticket_id"`
	IdempotencyKey string `json:"idempotency_key"`
}
