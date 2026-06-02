package event

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/adapters"
	"tickets/entities"
)

func (h Handler) BookDeadNationTicket(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("Booking ticket in Dead Nation")

	show, err := h.showsRepository.ShowByID(ctx, event.ShowID.String())
	if err != nil {
		return err
	}

	return h.deadNationClient.BookTicket(ctx, adapters.DeadNationBooking{
		BookingID:         event.BookingID,
		DeadNationEventID: show.DeadNationID,
		NumberOfTickets:   event.NumberOfTickets,
		CustomerEmail:     event.CustomerEmail,
	})
}
