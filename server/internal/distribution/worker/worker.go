package worker

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/distribution/task"
)

var ErrHeartbeatMismatch = errors.New("worker: heartbeat mismatch")

type Worker struct {
	// ID is the unique ID of the worker.
	ID uuid.UUID

	// Cores is the number of cores the worker contributes to the server.
	Cores uint

	// JoinedAt is the time the worker joined.
	JoinedAt time.Time

	// Tasks is the pool of tasks assigned to the worker.
	Tasks *task.Pool

	// removeCh is a channel that is closed when the worker is removed.
	removeCh chan struct{}
	removed  bool
}

// New creates a new worker.
func New(cores uint) *Worker {
	return &Worker{
		ID:       uuid.New(),
		Cores:    cores,
		JoinedAt: time.Now(),
		Tasks:    task.NewPool(),
		removeCh: make(chan struct{}, 1),
	}
}

// WorkerID returns the unique ID of the worker.
func (w *Worker) WorkerID() uuid.UUID {
	return w.ID
}

// WorkerCores returns the number of cores the worker contributes to the server.
func (w *Worker) WorkerCores() uint {
	return w.Cores
}

func (w *Worker) Remove() {
	if w.removed {
		return
	}
	w.removeCh <- struct{}{}
	close(w.removeCh)
	w.removed = true
}

func (w *Worker) NotifyRemoval() <-chan struct{} {
	return w.removeCh
}
