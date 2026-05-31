package event

import (
	"context"

	"tickets/entities"
)

type Handler struct {
	spreadsheetsAPI   SpreadsheetsAPI
	receiptsService   ReceiptsService
	ticketsRepository TicketsRepository
	printerService    PrinterService
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	ticketsRepository TicketsRepository,
	printerService PrinterService,
) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if ticketsRepository == nil {
		panic("missing ticketsRepository")
	}
	if printerService == nil {
		panic("missing printerService")
	}

	return Handler{
		spreadsheetsAPI:   spreadsheetsAPI,
		receiptsService:   receiptsService,
		ticketsRepository: ticketsRepository,
		printerService:    printerService,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type PrinterService interface {
	Print(ctx context.Context, ticket entities.Ticket) error
}
