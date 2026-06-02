package event

import (
	"context"
	"tickets/adapters"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	spreadsheetsAPI   SpreadsheetsAPI
	receiptsService   ReceiptsService
	filesAPI          FilesAPI
	ticketsRepository TicketsRepository
	showsRepository   ShowsRepository
	deadNationClient  DeadNationClient
	eventBus          *cqrs.EventBus
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	filesAPI FilesAPI,
	ticketsRepository TicketsRepository,
	showsRepository ShowsRepository,
	deadNationClient DeadNationClient,
	eventBus *cqrs.EventBus,
) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if filesAPI == nil {
		panic("missing filesAPI")
	}
	if ticketsRepository == nil {
		panic("missing ticketsRepository")
	}
	if showsRepository == nil {
		panic("missing showsRepository")
	}
	if deadNationClient == nil {
		panic("missing deadNationClient")
	}
	if eventBus == nil {
		panic("missing eventBus")
	}

	return Handler{
		spreadsheetsAPI:   spreadsheetsAPI,
		receiptsService:   receiptsService,
		filesAPI:          filesAPI,
		ticketsRepository: ticketsRepository,
		showsRepository:   showsRepository,
		deadNationClient:  deadNationClient,
		eventBus:          eventBus,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

type FilesAPI interface {
	UploadFile(ctx context.Context, fileID string, fileContent string) error
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type ShowsRepository interface {
	ShowByID(ctx context.Context, showID string) (entities.Show, error)
}

type DeadNationClient interface {
	BookTicket(ctx context.Context, booking adapters.DeadNationBooking) error
}
