package main

import (
	"context"
	"log/slog"
	"os"
	"tickets/worker"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	w := worker.NewWorker(spreadsheetsAPI, receiptsService)
	go w.Run(context.Background())

	err = service.New(
		w,
	).Run(context.Background())
	if err != nil {
		panic(err)
	}

}
