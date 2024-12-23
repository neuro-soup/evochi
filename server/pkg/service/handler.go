package service

import (
	"github.com/neuro-soup/evochi/server/internal/event"
	"github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
)

type Handler struct {
	workers *worker.Pool
	events  *event.Queue
}

var _ evochiv1connect.EvochiServiceHandler = (*Handler)(nil)

// NewHandler creates a new service handler with the given pool and queue.
func NewHandler(workers *worker.Pool, events *event.Queue) *Handler {
	return &Handler{workers: workers, events: events}
}
