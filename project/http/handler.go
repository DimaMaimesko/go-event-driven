package http

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Handler struct {
	publisher message.Publisher
	eventBus  cqrs.EventBus
}
