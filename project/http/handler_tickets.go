package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"tickets/worker"
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
		h.worker.Send(
			worker.Message{
				Task:     worker.TaskIssueReceipt,
				TicketID: ticket,
			},
			worker.Message{
				Task:     worker.TaskAppendToTracker,
				TicketID: ticket,
			},
		)
	}

	return c.NoContent(http.StatusOK)
}
