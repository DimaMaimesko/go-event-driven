package command

import (
	"context"
	"fmt"

	"tickets/entities"
)

func (h Handler) RefundTicket(ctx context.Context, ticketRefund *entities.RefundTicket) error {
	idempotencyKey := ticketRefund.Header.IdempotencyKey
	if idempotencyKey == "" {
		return fmt.Errorf("idempotency key is required")
	}

	err := h.receiptsServiceClient.VoidReceipt(ctx, entities.VoidReceipt{
		TicketID:       ticketRefund.TicketID,
		Reason:         "ticket refunded",
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to void receipt: %w", err)
	}

	return nil
}
