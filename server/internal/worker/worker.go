package worker

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
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

	// lastHeartbeat is the most recent heartbeat received from the worker.
	lastHeartbeat *heartbeat
}

// New creates a new worker.
func New(cores uint, attrs map[string][]byte) *Worker {
	return &Worker{
		ID:       uuid.New(),
		Cores:    cores,
		Attrs:    attrs,
		JoinedAt: time.Now(),
	}
}

// Healthy returns whether the worker has a healthy heartbeat.
func (w *Worker) Healthy(timeout time.Duration) bool {
	return w.lastHeartbeat == nil || w.lastHeartbeat.Ping() < timeout
}

// Heartbeat updates the worker's heartbeat information.
func (w *Worker) Heartbeat(seqID uint, sent time.Time) error {
	slog.Debug("worker heartbeat", "seq_id", seqID, "sent", sent, "worker", w.ID)

	if w.lastHeartbeat != nil && w.lastHeartbeat.SeqID != seqID-1 {
		return ErrHeartbeatMismatch
	}

	w.lastHeartbeat = &heartbeat{
		SeqID:      seqID,
		SentAt:     sent,
		ReceivedAt: time.Now(),
	}

	return nil
}
