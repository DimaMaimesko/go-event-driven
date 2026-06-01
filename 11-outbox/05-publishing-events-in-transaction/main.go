package main

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func PublishInTx(
	message *message.Message,
	tx *sqlx.Tx,
) error {
	// TODO: your code goes here
	return nil
}
