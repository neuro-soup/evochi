package worker

import (
	"time"

	"github.com/google/uuid"
)

type WorkerConfig struct {
	// Cores is the number of cores the worker contributes to the server.
	Cores uint

	// Device is the compute device the worker uses (e.g. "GPU", "CPU", ...)
	Device string

	// Attrs are custom attributes assigned by the worker.
	Attrs map[string]any
}

type Worker struct {
	// ID is the unique ID of the worker.
	ID uuid.UUID

	// Config is the configuration of the worker.
	Config WorkerConfig

	// JoinedAt is the time the worker joined.
	JoinedAt time.Time

	// Ping is the most recent travel time of heartbeats between worker and server.
	Ping time.Duration
}
