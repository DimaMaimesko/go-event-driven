package db_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tickets/db"
	"tickets/entities"
)

func TestTicketsRepository_IdempotentAdd(t *testing.T) {
	dbConn, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	require.NoError(t, err)
	defer dbConn.Close()

	err = db.InitializeDatabaseSchema(dbConn)
	require.NoError(t, err)

	repo := db.NewTicketsRepository(dbConn)

	ctx := context.Background()

	ticketID := uuid.NewString()
	ticket := entities.Ticket{
		TicketID: ticketID,
		Price: entities.Money{
			Amount:   "100.00",
			Currency: "USD",
		},
		CustomerEmail: "test@example.com",
	}

	// 1. Add the ticket for the first time
	err = repo.Add(ctx, ticket)
	require.NoError(t, err, "first addition should not fail")

	// 2. Add the ticket again - it should be idempotent
	err = repo.Add(ctx, ticket)
	require.NoError(t, err, "second addition should be idempotent and not fail")

	// 3. Verify that only one ticket was added
	tickets, err := repo.FindAll(ctx)
	require.NoError(t, err)

	count := 0
	for _, tkt := range tickets {
		if tkt.TicketID == ticketID {
			count++
		}
	}
	assert.Equal(t, 1, count, "there should be exactly one ticket with the given ID")
}
