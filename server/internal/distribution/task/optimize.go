package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var optimizeHist = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "optimize_duration_seconds",
	Help: "Optimize duration in seconds",
})

type Optimize struct {
	id uuid.UUID

	RequestedAt time.Time
	Timeout     time.Duration

	Epoch   uint
	Rewards []float64
}

var _ Task = (*Optimize)(nil)

func NewOptimize(epoch uint, rewards []float64, timeout time.Duration) *Optimize {
	return &Optimize{
		id:          uuid.New(),
		Epoch:       epoch,
		Rewards:     rewards,
		RequestedAt: time.Now(),
		Timeout:     timeout,
	}
}

func (i *Optimize) ID() uuid.UUID {
	return i.id
}

func (i *Optimize) Deadline() time.Time {
	return i.RequestedAt.Add(i.Timeout)
}

func (i *Optimize) Done() {
	optimizeHist.Observe(time.Since(i.RequestedAt).Seconds())
}
