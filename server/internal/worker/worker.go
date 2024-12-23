package worker

import (
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	// ID is the unique ID of the worker.
	ID uuid.UUID

	// Cores is the number of cores the worker contributes to the server.
	Cores uint

	// Attrs are custom attributes assigned by the worker.
	Attrs map[string][]byte

	// JoinedAt is the time the worker joined.
	JoinedAt time.Time

	// Ping is the most recent travel time of heartbeats between worker and server.
	Ping time.Duration

	// LastSeen is the time the worker was last seen.
	LastSeen time.Time
}

// NewWorker creates a new worker.
func NewWorker(cores uint, attrs map[string][]byte) *Worker {
	return &Worker{
		ID:       uuid.New(),
		Cores:    cores,
		Attrs:    attrs,
		JoinedAt: time.Now(),
	}
}
