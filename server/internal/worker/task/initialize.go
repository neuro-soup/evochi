package task

import (
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var initializeHist = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "initialize_duration_seconds",
	Help: "Initialize duration in seconds",
})

type Initialize struct {
	id uuid.UUID

	RequestedAt time.Time
	Timeout     time.Duration

	Epoch uint
}

var _ Task = (*Initialize)(nil)

func NewInitialize(epoch uint, timeout time.Duration) *Initialize {
	return &Initialize{
		id:          uuid.New(),
		Epoch:       epoch,
		RequestedAt: time.Now(),
		Timeout:     timeout,
	}
}

func Initializes(pool *Pool) []*Initialize {
	var inits []*Initialize
	for _, task := range pool.Tasks() {
		hb, ok := task.(*Initialize)
		if !ok {
			continue
		}
		inits = append(inits, hb)
	}
	return inits
}

func (i *Initialize) ID() uuid.UUID {
	return i.id
}

func (i *Initialize) Deadline() time.Time {
	return i.RequestedAt.Add(i.Timeout)
}

func (i *Initialize) Done() {
	initializeHist.Observe(time.Since(i.RequestedAt).Seconds())
}
