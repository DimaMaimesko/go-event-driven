package outbox

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
)

const outboxTopic = "events_to_forward"

func NewPublisherForDb(ctx context.Context, tx *sqlx.Tx) (message.Publisher, error) {
	logger := watermill.NewSlogLogger(log.FromContext(ctx))

	var publisher message.Publisher
	publisher, err := watermillSQL.NewPublisher(
		tx,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create publisher: %w", err)
	}

	publisher = log.CorrelationPublisherDecorator{Publisher: publisher}

	publisher = forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})

	publisher = log.CorrelationPublisherDecorator{Publisher: publisher}

	return publisher, nil
}
