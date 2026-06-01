package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	ticketsRepository TicketsRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:    eventBus,
		ticketsRepo: ticketsRepository,
	}

	//e.POST("/tickets-status", handler.PostTicketsStatus)
	e.POST("/tickets-status", handler.PostTicketsStatus, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Idempotency-Key") == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Idempotency-Key header is required")
			}
			return next(c)
		}
	})

	e.GET("/tickets", handler.GetTickets)

	return e
}
