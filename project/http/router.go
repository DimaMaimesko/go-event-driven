package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
