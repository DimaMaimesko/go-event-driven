package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/message/event"
)

func NewWatermillRouter(receiptsService event.ReceiptsService, spreadsheetsAPI event.SpreadsheetsAPI, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

	handler := event.NewHandler(spreadsheetsAPI, receiptsService)

	useMiddlewares(router, watermillLogger)

	ep, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(
					redisstream.SubscriberConfig{
						Client:        rdb,
						ConsumerGroup: params.HandlerName,
					},
					watermillLogger,
				)
			},
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
			Logger: watermillLogger,
		},
	)
	if err != nil {
		panic(err)
	}

	err = ep.AddHandlers(
		cqrs.NewEventHandler(
			"issue-receipt",
			handler.IssueReceipt,
		),
		cqrs.NewEventHandler(
			"append-to-tracker",
			handler.AppendToTracker,
		),
		cqrs.NewEventHandler(
			"cancel-ticket",
			handler.CancelTicket,
		),
	)
	if err != nil {
		panic(err)
	}

	return router
}
