package service

import (
	"context"
	"errors"
	stdHTTP "net/http"
	"os"
	"tickets/message"

	ticketsHttp "tickets/http"
	"tickets/worker"

	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type Service struct {
	echoRouter *echo.Echo
	publisher  watermillMessage.Publisher
}

func New(
	spreadsheetsAPI worker.SpreadsheetsAPI,
	receiptsService worker.ReceiptsService,
) Service {
	logger := watermill.NewSlogLogger(nil)
	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))

	publisher := message.NewRedisPublisher(redisClient, logger)

	message.NewSubscriber(
		redisClient,
		logger,
		"issue-receipt",
		"receipt",
		func(ctx context.Context, payload string) error {
			return receiptsService.IssueReceipt(ctx, payload)
		},
	)

	message.NewSubscriber(
		redisClient,
		logger,
		"append-to-tracker",
		"tracker",
		func(ctx context.Context, payload string) error {
			return spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{payload})
		},
	)

	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	return Service{
		echoRouter: echoRouter,
		publisher:  publisher,
	}
}

func (s Service) Run(ctx context.Context) error {

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
