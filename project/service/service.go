package service

import (
	"context"
	"errors"
	stdHTTP "net/http"
	"tickets/worker"

	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	worker *worker.Worker,
) Service {
	echoRouter := ticketsHttp.NewHttpRouter(*worker)

	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
