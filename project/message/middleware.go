package message

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/lithammer/shortuuid/v3"
)

func useMiddlewares(router *message.Router, watermillLogger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)

	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermillLogger,
	}.Middleware)

	router.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) (events []*message.Message, err error) {
			ctx := msg.Context()

			reqCorrelationID := msg.Metadata.Get("correlation_id")
			if reqCorrelationID == "" {
				reqCorrelationID = shortuuid.New()
			}

			ctx = log.ToContext(ctx, slog.With("correlation_id", reqCorrelationID))
			ctx = log.ContextWithCorrelationID(ctx, reqCorrelationID)

			msg.SetContext(ctx)

			return h(msg)
		}
	})

	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			logger := log.FromContext(msg.Context()).With(
				"message_id", msg.UUID,
				"payload", string(msg.Payload),
				"metadata", msg.Metadata,
				"handler", message.HandlerNameFromCtx(msg.Context()),
			)

			logger.Info("Handling a message")

			msgs, err := next(msg)
			if err != nil {
				logger.With("error", err).Error("Error while handling a message")
			}

			return msgs, err
		}
	})

	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			if msg.UUID == "5f810ce3-222b-4626-bc04-cbfb460c98c7" {
				logger := log.FromContext(msg.Context()).With(
					"message_id", msg.UUID,
					"payload", string(msg.Payload),
					"metadata", msg.Metadata,
					"handler", message.HandlerNameFromCtx(msg.Context()),
				)
				logger.Error("Error message", msg.UUID)
				return nil, nil
			}
			return next(msg)
		}
	})

	// Skip a known malformed message (invalid JSON payload published
	// to TicketBookingConfirmed). We ack it by returning no error.
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			if msg.UUID == "2beaf5bc-d5e4-4653-b075-2b36bbf28949" {
				log.FromContext(msg.Context()).With(
					"message_id", msg.UUID,
				).Info("Skipping known malformed message")
				return nil, nil
			}
			return next(msg)
		}
	})
}
