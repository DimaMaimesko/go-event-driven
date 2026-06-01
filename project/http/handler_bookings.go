package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"tickets/entities"
)

func (h Handler) BookTickets(c echo.Context) error {
	booking := entities.Booking{}

	if err := c.Bind(&booking); err != nil {
		return err
	}

	booking.BookingID = uuid.NewString()

	if err := h.bookingsRepository.BookTickets(c.Request().Context(), booking); err != nil {
		return fmt.Errorf("failed to add show: %w", err)
	}

	return c.JSON(http.StatusCreated, booking)
}
