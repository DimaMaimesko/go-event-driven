// This file contains tests that are executed to verify your solution.
// It's read-only, so all modifications will be ignored.
package db_test

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ticketsDb "tickets/db"
	"tickets/entities"
)

func TestBookingsRepository_AddBooking_seats_limit(t *testing.T) {
	ctx := context.Background()

	db := getDb()

	err := ticketsDb.InitializeDatabaseSchema(db)
	require.NoError(t, err)

	bookingsRepo := ticketsDb.NewBookingsRepository(db)
	showsRepo := ticketsDb.NewShowsRepository(db)

	t.Run("overbooking", func(t *testing.T) {
		showID := uuid.New()

		err := showsRepo.AddShow(ctx, entities.Show{
			ShowID:          showID,
			DeadNationID:    uuid.New(),
			NumberOfTickets: 2,
			StartTime:       time.Now().Add(time.Hour),
			Title:           "Example title",
			Venue:           "Exmaple vanue",
		})
		require.NoError(t, err)

		err = bookingsRepo.AddBooking(ctx, entities.Booking{
			BookingID:       uuid.New(),
			ShowID:          showID,
			NumberOfTickets: 2,
			CustomerEmail:   "foo@bar.com",
		})
		require.NoError(t, err)

		err = bookingsRepo.AddBooking(ctx, entities.Booking{
			BookingID:       uuid.New(),
			ShowID:          showID,
			NumberOfTickets: 2,
			CustomerEmail:   "foo@bar.com",
		})
		requireNotEnoughSeatsError(t, err)
	})

	t.Run("parallel_overbooking", func(t *testing.T) {
		showID := uuid.New()

		workersCount := 50
		workersErrs := make(chan error, workersCount)

		unlock := make(chan struct{})

		err := showsRepo.AddShow(ctx, entities.Show{
			ShowID:          showID,
			DeadNationID:    uuid.New(),
			NumberOfTickets: 2,
			StartTime:       time.Now().Add(time.Hour),
			Title:           "Example title",
			Venue:           "Exmaple vanue",
		})
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(workersCount)

		for range workersCount {
			go func() {
				defer wg.Done()

				// we are synchronizing goroutines to make sure that chance of overbooking is as high as possible
				<-unlock
				err = bookingsRepo.AddBooking(ctx, entities.Booking{
					BookingID:       uuid.New(),
					ShowID:          showID,
					NumberOfTickets: 2,
					CustomerEmail:   "foo@bar.com",
				})
				workersErrs <- err
			}()
		}
		close(unlock)

		wg.Wait()
		close(workersErrs)

		failedWorkers := 0
		succeededWorkers := 0
		errors := []error{}

		for err := range workersErrs {
			if err != nil {
				failedWorkers++
				errors = append(errors, err)
			} else {
				succeededWorkers++
			}
		}

		assert.Equal(t, 1, succeededWorkers)
		assert.Equal(t, workersCount-1, failedWorkers)

		if succeededWorkers == 0 {
			// all workers failed, let's print error
			for _, err := range errors {
				t.Log("error:", err)
			}
		}
	})
}

func requireNotEnoughSeatsError(t *testing.T, err error) {
	var echoErr *echo.HTTPError
	require.ErrorAs(t, err, &echoErr)

	require.Equal(t, http.StatusBadRequest, echoErr.Code)
	require.Equal(t, "not enough seats available", echoErr.Message)
}
