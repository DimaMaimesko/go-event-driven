package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SubscribeToMessages(
	db *sqlx.DB,
	topic string,
	logger watermill.LoggerAdapter,
) (<-chan *message.Message, error) {
	subscriber, err := watermillSQL.NewSubscriber(
		db,
		watermillSQL.SubscriberConfig{
			SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
			OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
		},
		logger,
	)
	if err != nil {
		return nil, nil
	}
	err = subscriber.SubscribeInitialize(topic)
	if err != nil {
		return nil, nil
	}
	message, err := subscriber.Subscribe(context.Background(), "ItemAddedToCart")

	return message, nil
}
