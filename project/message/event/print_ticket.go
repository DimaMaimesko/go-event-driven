package event

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) PrintTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Print ticket")

	return h.printerService.Print(
		ctx,
		entities.Ticket{
			TicketID:      event.TicketID,
			Price:         event.Price,
			CustomerEmail: event.CustomerEmail,
		},
	)
}
