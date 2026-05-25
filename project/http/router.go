package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(p message.Publisher) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		publisher: p,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
