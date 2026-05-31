package http

import (
	"tickets/db"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eventBus    *cqrs.EventBus
	ticketsRepo db.TicketsRepository
}
