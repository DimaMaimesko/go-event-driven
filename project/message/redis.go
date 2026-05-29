package message

import (
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
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
	pub = log.CorrelationPublisherDecorator{pub}

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
