package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"

	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter      *echo.Echo
	watermillRouter *watermillMessage.Router
}

func New(
	redisClient *redis.Client,
	spreadsheetsAPI message.SpreadsheetsAPI,
	receiptsService message.ReceiptsService,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)
	watermillRouter := watermillMessage.NewDefaultRouter(watermillLogger)

	message.NewHandlers(
		receiptsService,
		spreadsheetsAPI,
		redisClient,
		watermillLogger,
		watermillRouter,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		redisPublisher,
	)

	return Service{
		echoRouter,
		watermillRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	go func() {
		err := s.watermillRouter.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
