package v1

import (
	"github.com/neuro-soup/evochi/server/internal/epoch"
	"github.com/neuro-soup/evochi/server/internal/event"
	"github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
)

type Config struct {
	// JWTSecret is the secret used to sign JWT tokens.
	JWTSecret string

	// MaxWorkers is the maximum number of workers to start. If zero, there is
	// no limit.
	MaxWorkers uint

	// MaxEpochs is the maximum number of epochs to process. If zero, there is
	// no limit.
	MaxEpochs uint
}

type Handler struct {
	config Config

	workers *worker.Pool
	events  *event.Queue
	conns   *connections

	// epoch is the current epoch.
	epoch *epoch.Epoch
}

var _ evochiv1connect.EvochiServiceHandler = (*Handler)(nil)

// New creates a new handler.
func New(config Config, workers *worker.Pool, events *event.Queue) *Handler {
	return &Handler{
		config:  config,
		workers: workers,
		events:  events,
		conns:   newConnections(),
	}
}
