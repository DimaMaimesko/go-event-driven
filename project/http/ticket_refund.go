package http

import (
	"net/http"
	"tickets/entities"

	"github.com/labstack/echo/v4"
)

func (h Handler) TicketRefund(c echo.Context) error {
	ticketID := c.Param("ticket_id")

	cmd := entities.RefundTicket{
		Header:   entities.NewMessageHeader(),
		TicketID: ticketID,
	}

	err := h.commandBus.Send(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}
