package worker

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/worker/task"
)

var ErrHeartbeatMismatch = errors.New("worker: heartbeat mismatch")

type Worker struct {
	// ID is the unique ID of the worker.
	ID uuid.UUID

	// Cores is the number of cores the worker contributes to the server.
	Cores uint

	// Attrs are custom attributes assigned by the worker.
	Attrs map[string][]byte

	// JoinedAt is the time the worker joined.
	JoinedAt time.Time

	// Tasks is the pool of tasks assigned to the worker.
	Tasks *task.Pool

	removeCh chan struct{}
}

// New creates a new worker.
func New(cores uint, attrs map[string][]byte) *Worker {
	return &Worker{
		ID:       uuid.New(),
		Cores:    cores,
		Attrs:    attrs,
		JoinedAt: time.Now(),
		Tasks:    task.NewPool(),
		removeCh: make(chan struct{}),
	}
}

func (w *Worker) Remove() {
	w.removeCh <- struct{}{}
	close(w.removeCh) // TODO: validate this
}

func (w *Worker) Removes() <-chan struct{} {
	return w.removeCh
}
