package service

import (
	"context"
	"errors"
	"fmt"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"tickets/db"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"
)

const outboxTopic = "events_to_forward"

type Service struct {
	db              *sqlx.DB
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	spreadsheetsAPI event.SpreadsheetsAPI,
	receiptsService event.ReceiptsService,
	filesAPI event.FilesAPI,
) Service {
	ticketsRepo := db.NewTicketsRepository(dbConn)
	showsRepo := db.NewShowsRepository(dbConn)
	bookingsRepository := db.NewBookingsRepository(dbConn)

	watermillLogger := watermill.NewSlogLogger(log.FromContext(context.Background()))

	redisPublisher := message.NewRedisPublisher(redisClient, watermillLogger)

	eventBus := event.NewBus(redisPublisher)

	eventsHandler := event.NewHandler(
		spreadsheetsAPI,
		receiptsService,
		filesAPI,
		ticketsRepo,
		eventBus,
	)
	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		eventProcessorConfig,
		eventsHandler,
		watermillLogger,
	)

	// SQL subscriber reading enveloped messages from the outbox table.
	// InitializeSchema=true makes sure the watermill_events_to_forward table is created automatically.
	sqlSubscriber, err := watermillSQL.NewSubscriber(
		dbConn,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:    watermillSQL.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
			InitializeSchema: true,
		},
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}

	// The Forwarder reads enveloped messages from the SQL subscriber and re-publishes
	// them to their destination topic using the existing Redis publisher.
	// Passing the existing Router via Config.Router registers the Forwarder's handler
	// with it, so the Forwarder runs as part of the Router's lifecycle.
	_, err = forwarder.NewForwarder(
		sqlSubscriber,
		redisPublisher,
		watermillLogger,
		forwarder.Config{
			ForwarderTopic: outboxTopic,
			Router:         watermillRouter,
		},
	)
	if err != nil {
		panic(err)
	}

	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
		ticketsRepo,
		showsRepo,
		bookingsRepository,
	)

	return Service{
		dbConn,
		watermillRouter,
		echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDatabaseSchema(s.db); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
