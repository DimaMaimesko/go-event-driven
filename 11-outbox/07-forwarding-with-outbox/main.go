package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	forwarder2 "github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {
	postgresSub, err := watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
			OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	err = postgresSub.SubscribeInitialize(outboxTopic)
	if err != nil {
		return err
	}

	redisPub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		return err
	}

	forwarder, err := forwarder2.NewForwarder(
		postgresSub,
		redisPub,
		logger,
		forwarder2.Config{
			ForwarderTopic: outboxTopic,
		},
	)
	if err != nil {
		return err
	}

	go func() {
		err := forwarder.Run(context.Background())
		if err != nil {
			logger.Error("forwarder failed", err, nil)
		}
	}()

	<-forwarder.Running()

	return nil
}
