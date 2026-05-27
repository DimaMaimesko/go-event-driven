package events

import (
	"tickets/entities"
	"time"

	"github.com/google/uuid"
)

type TicketBookingConfirmed struct {
	Header        MessageHeader  `json:"header"`
	TicketID      string         `json:"ticket_id"`
	CustomerEmail string         `json:"customer_email"`
	Price         entities.Money `json:"price"`
}

type MessageHeader struct {
	ID          string    `json:"id"`
	PublishedAt time.Time `json:"published_at"`
}

func NewMessageHeader() MessageHeader {
	return MessageHeader{
		ID:          uuid.NewString(),
		PublishedAt: time.Now().UTC(),
	}
}
