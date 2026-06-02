package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/jmoiron/sqlx"

	"tickets/entities"
)

type BookingsRepository struct {
	db *sqlx.DB
}

func NewBookingsRepository(db *sqlx.DB) BookingsRepository {
	if db == nil {
		panic("nil db")
	}

	return BookingsRepository{db: db}
}

func (b BookingsRepository) AddBooking(ctx context.Context, booking entities.Booking) (err error) {
	return updateInTx(
		ctx,
		b.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.NamedExecContext(ctx, `
				INSERT INTO 
				    bookings (booking_id, show_id, number_of_tickets, customer_email) 
				VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
				`, booking)
			if err != nil {
				return fmt.Errorf("could not add booking: %w", err)
			}

			outboxPublisher, err := outbox.NewPublisherForDb(ctx, tx)
			if err != nil {
				return fmt.Errorf("could not create event bus publisher: %w", err)
			}

			bus := event.NewBus(outboxPublisher)

			err = bus.Publish(ctx, entities.BookingMade{
				Header:          entities.NewMessageHeader(),
				NumberOfTickets: booking.NumberOfTickets,
				BookingID:       booking.BookingID.String(),
				CustomerEmail:   booking.CustomerEmail,
				ShowID:          booking.ShowID.String(),
			})
			if err != nil {
				return fmt.Errorf("could not publish BookingMade event: %w", err)
			}

			return nil
		},
	)
}

func updateInTx(
	ctx context.Context,
	db *sqlx.DB,
	isolation sql.IsolationLevel,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
			return
		}

		err = tx.Commit()
	}()

	return fn(ctx, tx)
}
