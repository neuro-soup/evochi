package event

import (
	"log/slog"
	"reflect"
	"sync"

	"github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var eventsPushed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "event_queue_events_pushed",
	Help: "The total number of events pushed to the queue.",
})

type Queue struct {
	mu    *sync.Mutex
	chans map[*worker.Worker]chan Event
}

// NewQueue creates a new event queue.
func NewQueue() *Queue {
	return &Queue{
		mu:    new(sync.Mutex),
		chans: make(map[*worker.Worker]chan Event),
	}
}

// Push pushes event to the worker's channel.
func (q *Queue) Push(w *worker.Worker, evt Event) {
	q.mu.Lock()
	defer q.mu.Unlock()

	slog.Debug("pushing event to worker", "worker", w.ID, "event", reflect.TypeOf(evt))

	if q.chans[w] == nil {
		q.chans[w] = make(chan Event)
	}
	q.chans[w] <- evt

	eventsPushed.Inc()
}

// Pull returns a channel for the worker to pull events from.
func (q *Queue) Pull(w *worker.Worker) <-chan Event {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.chans[w] == nil {
		q.chans[w] = make(chan Event)
	}
	return q.chans[w]
}

// Close closes the worker's channel.
func (q *Queue) Close(w *worker.Worker) {
	q.mu.Lock()
	defer q.mu.Unlock()

	close(q.chans[w])
	delete(q.chans, w)
}
