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
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			err = h.publisher.Publish("TicketBookingConfirmed", msg)
			if err != nil {
				return err
			}

			//publish malformed message
			eventMalformed := entities.TicketBookingConfirmed{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: "dededed",
				Price:         ticket.Price,
			}

			Malformed, err := json.Marshal(eventMalformed)
			if err != nil {
				return err
			}

			msg = message.NewMessage("2beaf5bc-d5e4-4653-b075-2b36bbf28949", Malformed)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			msg.Metadata.Set("type", "TicketBookingConfirmed")

			err = h.publisher.Publish("TicketBookingConfirmed", msg)
			if err != nil {
				return err
			}

		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			err = h.publisher.Publish("TicketBookingCanceled", msg)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}
