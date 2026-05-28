package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	ticketsHttp "tickets/http"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tickets/adapters"
	"tickets/entities"
	"tickets/message"
	"tickets/service"
)

func TestComponent(t *testing.T) {
	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spreadsheetsAPI := &adapters.SpreadsheetsAPIStub{}
	receiptsService := &adapters.ReceiptsServiceStub{}

	go func() {
		svc := service.New(
			redisClient,
			spreadsheetsAPI,
			receiptsService,
		)
		assert.NoError(t, svc.Run(ctx))
	}()

	waitForHttpServer(t)

	ticketID := shortuuid.New()
	price := entities.Money{
		Amount:   "50.30",
		Currency: "GBP",
	}

	sendTicketsStatus(t, ticketsHttp.TicketsStatusRequest{
		Tickets: []ticketsHttp.TicketStatusRequest{
			{
				TicketID:      ticketID,
				Status:        "confirmed",
				Price:         price,
				CustomerEmail: "email@example.com",
			},
		},
	})

	assertReceiptForTicketIssued(t, receiptsService, ticketID, price)
	assertRowToSheetAdded(t, spreadsheetsAPI, "tickets-to-print", ticketID)
}

func sendTicketsStatus(t *testing.T, req ticketsHttp.TicketsStatusRequest) {
	t.Helper()

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	correlationID := shortuuid.New()

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("Correlation-ID", correlationID)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func assertReceiptForTicketIssued(
	t *testing.T,
	receiptsService *adapters.ReceiptsServiceStub,
	ticketID string,
	expectedPrice entities.Money,
) {
	t.Helper()

	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			issuedReceipts := len(receiptsService.IssuedReceipts)
			assert.Greater(collectT, issuedReceipts, 0, "no receipts issued")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	var receipt entities.IssueReceiptRequest
	var found bool
	for _, r := range receiptsService.IssuedReceipts {
		if r.TicketID == ticketID {
			receipt = r
			found = true
			break
		}
	}

	require.True(t, found, "receipt for ticket %s not found", ticketID)
	assert.Equal(t, expectedPrice, receipt.Price)
}

func assertRowToSheetAdded(
	t *testing.T,
	spreadsheetsAPI *adapters.SpreadsheetsAPIStub,
	sheetName string,
	ticketID string,
) {
	t.Helper()

	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			rows, ok := spreadsheetsAPI.Rows[sheetName]
			if !assert.True(collectT, ok, "sheet %s not found", sheetName) {
				return
			}

			var found bool
			for _, row := range rows {
				if len(row) > 0 && row[0] == ticketID {
					found = true
					break
				}
			}
			assert.True(collectT, found, "row for ticket %s not found in sheet %s", ticketID, sheetName)
		},
		10*time.Second,
		100*time.Millisecond,
	)
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}
