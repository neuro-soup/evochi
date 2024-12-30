package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var shareStateHist = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "share_state_duration_seconds",
	Help: "Share state duration in seconds",
})

type ShareState struct {
	id uuid.UUID

	Epoch uint

	RequestedAt time.Time
	Timeout     time.Duration
}

var _ Task = (*ShareState)(nil)

func NewShareState(epoch uint, timeout time.Duration) *ShareState {
	return &ShareState{
		id:          uuid.New(),
		Epoch:       epoch,
		RequestedAt: time.Now(),
		Timeout:     timeout,
	}
}

func (e *ShareState) ID() uuid.UUID {
	return e.id
}

func (e *ShareState) Deadline() time.Time {
	return e.RequestedAt.Add(e.Timeout)
}

func (i *ShareState) Done() {
	shareStateHist.Observe(time.Since(i.RequestedAt).Seconds())
}
