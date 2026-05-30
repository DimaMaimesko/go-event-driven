package event

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) CancelTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Removing cancelled ticket")

	return h.ticketsRepository.Delete(ctx, event.TicketID)
}
