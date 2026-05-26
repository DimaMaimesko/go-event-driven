package http

import (
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
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
		msg := message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))

		err = h.publisher.Publish("issue-receipt", msg)
		if err != nil {
			return err
		}

		msg = message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))

		err = h.publisher.Publish("append-to-tracker", msg)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
