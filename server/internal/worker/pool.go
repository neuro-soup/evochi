package worker

import (
	"fmt"
	"iter"
	"sync"

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
)

type Pool struct {
	mu   *sync.RWMutex
	pool map[uuid.UUID]*Worker
}

// NewPool creates a new empty worker pool.
func NewPool() *Pool {
	return &Pool{
		mu: new(sync.RWMutex),
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
}

// Remove removes the worker from the pool.
func (p *Pool) Remove(id uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.pool, id)

	workersRemoved.Inc()
}

// Iter iterates over the workers in the pool.
func (p *Pool) Iter() iter.Seq[*Worker] {
	return func(yield func(*Worker) bool) {
		p.mu.RLock()
		defer p.mu.RUnlock()

		for _, w := range p.pool {
			if !yield(w) {
				return
			}
		}
	}
}
