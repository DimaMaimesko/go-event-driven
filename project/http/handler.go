package http

import (
	"tickets/worker"
)

type Handler struct {
	worker *worker.Worker
}
