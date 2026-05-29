package main

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewEventProcessor(
	router *message.Router,
	rdb *redis.Client,
	marshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
) (*cqrs.EventProcessor, error) {
	panic("implement me!")
}
