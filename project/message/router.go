package message

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/entities"
	"tickets/entities/events"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

func NewWatermillRouter(
	receiptsService ReceiptsService,
	spreadsheetsAPI SpreadsheetsAPI,
	rdb *redis.Client,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

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

	router.AddConsumerHandler(
		"issue_receipt",
		"TicketBookingConfirmed",
		issueReceiptSub,
		func(msg *message.Message) error {
			ctx := msg.Context()

			var event events.TicketBookingConfirmed
			if err := json.Unmarshal(msg.Payload, &event); err != nil {
				return err
			}

			slog.Info("Issuing receipt")

			request := entities.IssueReceiptRequest{
				TicketID: event.TicketID,
				Price:    event.Price,
			}

			if err := receiptsService.IssueReceipt(ctx, request); err != nil {
				return fmt.Errorf("failed to issue receipt: %w", err)
			}

			return nil
		},
	)

	router.AddConsumerHandler(
		"append_to_tracker",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		func(msg *message.Message) error {
			ctx := msg.Context()

			var event events.TicketBookingConfirmed
			if err := json.Unmarshal(msg.Payload, &event); err != nil {
				return err
			}

			slog.Info("Appending ticket to the tracker")

			return spreadsheetsAPI.AppendRow(
				ctx,
				"tickets-to-print",
				[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
			)
		},
	)

	return router
}
