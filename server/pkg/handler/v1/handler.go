package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/neuro-soup/evochi/server/internal/distribution/task"
	"github.com/neuro-soup/evochi/server/internal/distribution/worker"
	"github.com/neuro-soup/evochi/server/internal/training/epoch"
	"github.com/neuro-soup/evochi/server/internal/training/eval"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
)

type Config struct {
	// JWTPrivateKey is the private key to use for JWT tokens.
	JWTSecret string

	// MaxWorkers is the maximum number of workers to start. If zero, there is
	// no limit.
	MaxWorkers uint

	// WorkerTimeout is the timeout for a worker. Defaults to 1 minute.
	WorkerTimeout time.Duration

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

func (h *Handler) reset() {
	slog.Debug("resetting epoch")

	h.epoch = nil
	for _, w := range h.workers.Workers() {
		w.Remove()
	}
}

// eval tries to assign slices to the worker.
func (h *Handler) eval(w *worker.Worker) {
	// try to assign slices to the worker
	assigned := h.epoch.Assign(w)
	if len(assigned) == 0 {
		// no tasks to assign, worker is idle
		slog.Debug("no slices assigned, worker is idle", "worker", w.ID)

		if h.finished() {
			// all workers are idle, optimize
			h.optimize()
		}
		return
	}

	// add task to worker, which will be picked up by the event-loop
	t := task.NewEvaluate(h.epoch.Number, assigned, h.cfg.WorkerTimeout)
	w.Tasks.Add(t)

	slog.Debug("assigned slices to worker", "worker", w.ID, "slices", assigned)
}

func (h *Handler) finished() bool {
	for _, w := range h.workers.Workers() {
		if !w.Tasks.Idle() {
			fmt.Println("worker not idleeeeeeeeeeeeeeeeeee")
			b, _ := json.MarshalIndent(w.Tasks.Tasks(), "", "  ")
			fmt.Println(string(b))
			fmt.Println(w.Tasks.Tasks())
			return false
		}
	}
	return true
}

func (h *Handler) optimize() {
	slog.Debug("all workers finished, optimizing")
	rewards := h.epoch.Rewards()
	for _, w := range h.workers.Workers() {
		w.Tasks.Add(task.NewOptimize(h.epoch.Number, rewards, h.cfg.WorkerTimeout))
	}
}

func (h *Handler) requestState() {
	t := h.workers.Trusted(nil)
	if t == nil {
		slog.Error("no trusted worker available, cannot request state")
		return
	}

	t.Tasks.Add(task.NewShareState(h.epoch.Number, h.cfg.WorkerTimeout))
}

func (h *Handler) nextEpoch(state eval.State) {
	slog.Debug("starting next epoch", "epoch", h.epoch.Number+1)

	if h.cfg.MaxEpochs > 0 && h.epoch.Number+1 > h.cfg.MaxEpochs {
		slog.Info("reached maximum number of epochs", "max", h.cfg.MaxEpochs, "current", h.epoch.Number+1)
		// TODO:: maybe send a message to all workers to stop
		return
	}

	// go to next epoch
	h.epoch = epoch.New(h.epoch.Number+1, h.cfg.PopulationSize, state)

	// start evaluating next epoch
	for _, w := range h.workers.Workers() {
		h.eval(w)
	}
}
