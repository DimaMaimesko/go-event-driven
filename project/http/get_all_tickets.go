package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetTickets(c echo.Context) error {
	tickets, err := h.ticketsRepo.FindAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tickets)
}
