package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/training/eval"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var evaluateHist = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "evaluate_duration_seconds",
	Help: "Evaluate duration in seconds",
})

type Evaluate struct {
	id uuid.UUID

	RequestedAt time.Time
	Timeout     time.Duration

	Epoch  uint
	Slices []eval.Slice
}

var _ Task = (*Initialize)(nil)

func NewEvaluate(epoch uint, slices []eval.Slice, timeout time.Duration) *Evaluate {
	return &Evaluate{
		id:          uuid.New(),
		Epoch:       epoch,
		Slices:      slices,
		RequestedAt: time.Now(),
		Timeout:     timeout,
	}
}

func (e *Evaluate) ID() uuid.UUID {
	return e.id
}

func (e *Evaluate) Deadline() time.Time {
	return e.RequestedAt.Add(e.Timeout)
}

func (i *Evaluate) Done() {
	evaluateHist.Observe(time.Since(i.RequestedAt).Seconds())
}
