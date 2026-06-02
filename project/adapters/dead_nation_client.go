package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/dead_nation"
	"github.com/google/uuid"
)

type DeadNationClient struct {
	clients *clients.Clients
}

func NewDeadNationClient(clients *clients.Clients) *DeadNationClient {
	if clients == nil {
		panic("NewDeadNationClient: clients is nil")
	}

	return &DeadNationClient{clients: clients}
}

type DeadNationBooking struct {
	BookingID         uuid.UUID
	DeadNationEventID uuid.UUID
	NumberOfTickets   int
	CustomerEmail     string
}

func (c DeadNationClient) BookTicket(ctx context.Context, booking DeadNationBooking) error {
	resp, err := c.clients.DeadNation.PostTicketBookingWithResponse(
		ctx,
		dead_nation.PostTicketBookingRequest{
			BookingId:       booking.BookingID,
			EventId:         booking.DeadNationEventID,
			NumberOfTickets: booking.NumberOfTickets,
			// Dead Nation calls it "CustomerAddress", but for us it's the customer email.
			CustomerAddress: booking.CustomerEmail,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to post ticket booking to dead nation: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(
			"unexpected status code for POST dead-nation/ticket-booking: %d",
			resp.StatusCode(),
		)
	}

	return nil
}
