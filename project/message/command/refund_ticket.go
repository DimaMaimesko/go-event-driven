package command

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) RefundTicket(ctx context.Context, command *entities.RefundTicket) error {
	log.FromContext(ctx).Info("Refunding ticket")

	err := h.receiptsService.VoidReceipt(ctx, entities.VoidReceipt{
		TicketID:       command.TicketID,
		IdempotencyKey: command.Header.IdempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("failed to void receipt: %w", err)
	}

	return nil
}
