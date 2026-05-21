package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		err = h.receiptsService.IssueReceipt(c.Request().Context(), ticket)
		if err != nil {
			return err
		}

		err = h.spreadsheetsAPI.AppendRow(c.Request().Context(), "tickets-to-print", []string{ticket})
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
