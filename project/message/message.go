package message

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewRedisPublisher(rdb *redis.Client, watermillLogger watermill.LoggerAdapter) message.Publisher {
	var pub message.Publisher
	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	return pub
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
		// Go 1.25 made GOMAXPROCS container-aware, which lowers the go-redis default pool size
		// (10 x GOMAXPROCS) in small containers. With many consumer groups each holding
		// blocking XREADGROUP connections, publishers can starve. Raise the limit to be safe.
		PoolSize: 200,
	})
}

func NewSubscriber(
	rdb *redis.Client,
	watermillLogger watermill.LoggerAdapter,
	topic string,
	consumerGroup string,
	handler func(ctx context.Context, payload string) error,
) {
	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: consumerGroup,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	messages, err := sub.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range messages {
			if err := handler(msg.Context(), string(msg.Payload)); err != nil {
				slog.With(
					"error", err,
					"message_uuid", msg.UUID,
					"topic", topic,
					"consumer_group", consumerGroup,
				).Error("failed to handle message")
				msg.Nack()
				continue
			}
			msg.Ack()
		}
	}()
}
