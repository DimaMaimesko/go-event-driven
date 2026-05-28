package message

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/entities"
	"tickets/message/event"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

func NewWatermillRouter(receiptsService event.ReceiptsService, spreadsheetsAPI event.SpreadsheetsAPI, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

	handler := event.NewHandler(spreadsheetsAPI, receiptsService)

	useMiddlewares(router, watermillLogger)

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	cancelTicketSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "cancel-ticket",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	router.AddConsumerHandler(
		"issue_receipt",
		"TicketBookingConfirmed",
		issueReceiptSub,
		func(msg *message.Message) error {
			// Fixing a malformed JSON message
			// TODO: Remove once fixed
			if string(msg.UUID) == brokenMessageID {
				return nil
			}

			// Fixing an incorrect message type
			// TODO: Remove once fixed
			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}

			var event entities.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}

			return handler.IssueReceipt(msg.Context(), event)
		},
	)

	router.AddConsumerHandler(
		"append_to_tracker",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		func(msg *message.Message) error {
			// Fixing a malformed JSON message
			// TODO: Remove once fixed
			if string(msg.UUID) == brokenMessageID {
				return nil
			}

			// Fixing an incorrect message type
			// TODO: Remove once fixed
			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}

			var event entities.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}

			return handler.AppendToTracker(msg.Context(), event)
		},
	)

	router.AddConsumerHandler(
		"cancel_ticket",
		"TicketBookingCanceled",
		cancelTicketSub,
		func(msg *message.Message) error {
			// Fixing an incorrect message type
			// TODO: Remove once fixed
			if msg.Metadata.Get("type") != "TicketBookingCanceled" {
				return nil
			}

			var event entities.TicketBookingCanceled
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}
			return handler.CancelTicket(msg.Context(), event)
		},
	)

	return router
}
