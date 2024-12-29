package v1

import (
	"time"

	"github.com/neuro-soup/evochi/server/internal/distribution/worker"
	"github.com/neuro-soup/evochi/server/internal/training/epoch"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
)

type Config struct {
	// JWTSecret is the secret used to sign JWT tokens.
	JWTSecret string

	// MaxWorkers is the maximum number of workers to start. If zero, there is
	// no limit.
	MaxWorkers uint

	// WorkerTimeout is the timeout for a worker. Defaults to 1 minute.
	WorkerTimeout time.Duration `split_words:"true" default:"1m"`

	// MaxEpochs is the maximum number of epochs to process. If zero, there is
	// no limit.
	MaxEpochs uint

	// PopulationSize is the size of the population to use.
	PopulationSize uint

	// Attrs are the custom attributes to use.
	Attrs map[string][]byte
}

type Handler struct {
	cfg Config

	workers *worker.Pool

	// epoch is the current epoch.
	epoch *epoch.Epoch
}

var _ evochiv1connect.EvochiServiceHandler = (*Handler)(nil)

// New creates a new handler.
func New(config Config, workers *worker.Pool) *Handler {
	return &Handler{
		cfg:     config,
		workers: workers,
	}
}
