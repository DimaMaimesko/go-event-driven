package main

import (
	"context"
)

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request IssueReceiptRequest) error
}

type ReceiptsServiceStub struct {
	// todo: implement me
}
