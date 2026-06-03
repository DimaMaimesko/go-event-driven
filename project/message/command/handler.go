package command

import (
	"context"

	"tickets/entities"
)

type Handler struct {
	receiptsService ReceiptsService
}

func NewHandler(
	receiptsService ReceiptsService,
) Handler {
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		receiptsService: receiptsService,
	}
}

type ReceiptsService interface {
	VoidReceipt(ctx context.Context, request entities.VoidReceipt) error
}
