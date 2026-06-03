package command

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	eventBus *cqrs.EventBus

	receiptsServiceClient ReceiptsService
}

func NewHandler(eventBus *cqrs.EventBus, receiptsServiceClient ReceiptsService) Handler {
	if eventBus == nil {
		panic("eventBus is required")
	}
	if receiptsServiceClient == nil {
		panic("receiptsServiceClient is required")
	}

	handler := Handler{
		eventBus:              eventBus,
		receiptsServiceClient: receiptsServiceClient,
	}

	return handler
}

type ReceiptsService interface {
	VoidReceipt(ctx context.Context, request entities.VoidReceipt) error
}
