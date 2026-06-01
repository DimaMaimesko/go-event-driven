package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"tickets/entities"
)

type ShowRequest struct {
	DeadNationID    string `json:"dead_nation_id"`
	NumberOfTickets int    `json:"number_of_tickets"`
	StartTime       string `json:"start_time"`
	Title           string `json:"title"`
	Venue           string `json:"venue"`
}

type ShowResponse struct {
	ShowID string `json:"show_id"`
}

func (h Handler) Shows(c echo.Context) error {
	var request ShowRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	showID := uuid.NewString()

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		return fmt.Errorf("failed to parse start_time: %w", err)
	}

	show := entities.Show{
		ShowID:          showID,
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       startTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}

	err = h.showsRepo.AddShow(c.Request().Context(), show)
	if err != nil {
		return fmt.Errorf("failed to add show: %w", err)
	}

	return c.JSON(http.StatusCreated, ShowResponse{ShowID: showID})
}
