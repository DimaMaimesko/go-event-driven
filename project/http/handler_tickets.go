package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"

	"tickets/entities"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			msg := message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))

			err = h.publisher.Publish("issue-receipt", msg)
			if err != nil {
				return err
			}

			appendToTrackerPayload := entities.AppendToTrackerPayload{
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}
			payload, err := json.Marshal(appendToTrackerPayload)
			if err != nil {
				return err
			}
			msg = message.NewMessage(watermill.NewUUID(), payload)
			err = h.publisher.Publish("append-to-tracker", msg)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}
