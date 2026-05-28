package message

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/entities"
	"tickets/message/event"
)

// skipIfWrongType returns true if the message's "type" metadata is set and
// does not match the expected event type. Such messages are skipped and acked.
func skipIfWrongType(msg *message.Message, expectedType string) bool {
	msgType := msg.Metadata.Get("type")
	if msgType == "" {
		// No type metadata — accept as-is (backward compatible).
		return false
	}
	if msgType != expectedType {
		log.FromContext(msg.Context()).With(
			"message_id", msg.UUID,
			"expected_type", expectedType,
			"actual_type", msgType,
		).Info("Skipping message with mismatched type metadata")
		return true
	}
	return false
}

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
			if skipIfWrongType(msg, "TicketBookingConfirmed") {
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
			if skipIfWrongType(msg, "TicketBookingConfirmed") {
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
			if skipIfWrongType(msg, "TicketBookingCanceled") {
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
