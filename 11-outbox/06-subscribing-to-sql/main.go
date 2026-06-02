package main

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SubscribeToMessages(
	db *sqlx.DB,
	topic string,
	logger watermill.LoggerAdapter,
) (<-chan *message.Message, error) {
	// TODO: your code goes here
	return nil, nil
}
