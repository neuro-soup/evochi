package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var heartbeatHist = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "heartbeat_duration_seconds",
	Help: "Heartbeat duration in seconds",
})

type Heartbeat struct {
	id uuid.UUID

	SeqID uint

	RequestedAt time.Time
	Timeout     time.Duration
}

var _ Task = (*Heartbeat)(nil)

func NewHeartbeat(seqID uint, timeout time.Duration) *Heartbeat {
	return &Heartbeat{
		id:          uuid.New(),
		SeqID:       seqID,
		RequestedAt: time.Now(),
		Timeout:     timeout,
	}
}

func Heartbeats(pool *Pool) []*Heartbeat {
	var heartbeats []*Heartbeat
	for _, task := range pool.Tasks() {
		hb, ok := task.(*Heartbeat)
		if !ok {
			continue
		}
		heartbeats = append(heartbeats, hb)
	}
	return heartbeats
}

func (h *Heartbeat) ID() uuid.UUID {
	return h.id
}

func (h *Heartbeat) Deadline() time.Time {
	return h.RequestedAt.Add(h.Timeout)
}

func (h *Heartbeat) Done() {
	heartbeatHist.Observe(time.Since(h.RequestedAt).Seconds())
}
