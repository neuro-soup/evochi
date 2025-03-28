package worker

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	workersAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "worker_pool_workers_added",
		Help: "The total number of workers added to the pool.",
	})
	workersRemoved = promauto.NewCounter(prometheus.CounterOpts{
		Name: "worker_pool_workers_removed",
		Help: "The total number of workers removed from the pool.",
	})
	workers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "worker_pool_workers",
		Help: "The total number of workers in the pool.",
	})
)

type Pool struct {
	mu   *sync.RWMutex
	pool map[uuid.UUID]*Worker
}

// NewPool creates a new empty worker pool.
func NewPool() *Pool {
	return &Pool{
		mu:   new(sync.RWMutex),
		pool: make(map[uuid.UUID]*Worker),
	}
}

// String returns a string representation of the pool for debugging.
func (p *Pool) String() string {
	return fmt.Sprintf("Pool(Len=%d)", p.Len())
}

// Len returns the number of workers in the pool.
func (p *Pool) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.pool)
}

// Get returns the worker with the given ID, if it exists, or nil otherwise.
func (p *Pool) Get(id uuid.UUID) *Worker {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.pool[id]
}

// Add adds the worker to the pool.
func (p *Pool) Add(w *Worker) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pool[w.ID] = w

	workersAdded.Inc()
	workers.Set(float64(len(p.pool)))
}

// Remove removes the worker from the pool.
func (p *Pool) Remove(w *Worker) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.remove(w)
}

// remove removes the worker from the pool without locking.
func (p *Pool) remove(w *Worker) {
	delete(p.pool, w.ID)
	workersRemoved.Inc()

	w.Remove()
	workers.Set(float64(len(p.pool)))
}

func (p *Pool) Workers() []*Worker {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var workers []*Worker
	for _, w := range p.pool {
		workers = append(workers, w)
	}

	return workers
}

// Watch watches the pool for unproductive workers and removes them.
func (p *Pool) Watch(sleep time.Duration) {
	for {
		p.mu.Lock()
		for _, w := range p.pool {
			if w.Tasks.Punctual() {
				continue
			}

			slog.Info("removing unproductive worker", "worker", w.ID)
			p.remove(w)
		}
		p.mu.Unlock()

		time.Sleep(sleep)
	}
}

// Trusted returns a random trustworthy worker that passes the filter.
func (p *Pool) Trusted(filter func(w *Worker) bool) *Worker {
	var pool []*Worker
	for _, w := range p.pool {
		if (filter == nil || filter(w)) && w.Tasks.Punctual() {
			pool = append(pool, w)
		}
	}

	// no workers left
	if len(pool) == 0 {
		return nil
	}

	return pool[rand.Intn(len(pool))]
}

func (p *Pool) Idling() []*Worker {
	var idling []*Worker
	for _, w := range p.pool {
		if !w.Tasks.Idle() {
			continue
		}
		idling = append(idling, w)
	}
	return idling
}
