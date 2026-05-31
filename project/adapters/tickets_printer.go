package adapters

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
)

type PrinterClient struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewPrinterClient(clients *clients.Clients) *PrinterClient {
	if clients == nil {
		panic("NewReceiptsServiceClient: clients is nil")
	}

	return &PrinterClient{clients: clients}
}

func (c PrinterClient) Print(ctx context.Context, ticket entities.Ticket) error {
	ticketHTML := fmt.Sprintf(`
		<html>
			<head>
				<title>Ticket %s</title>
			</head>
			<body>
				<h1>Ticket ID: %s</h1>
				<p>Price: %s %s</p>
				<p>Customer Email: %s</p>
			</body>
		</html>
	`, ticket.TicketID, ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)

	resp, err := c.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(
		ctx,
		ticket.TicketID+"-ticket.html",
		ticketHTML,
	)
	if err != nil {
		return fmt.Errorf("failed to save ticket file: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("unexpected status code for PUT files: %d", resp.StatusCode())
	}

	return nil
}
