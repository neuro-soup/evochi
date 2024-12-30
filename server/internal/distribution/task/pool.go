package task

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	mu   *sync.Mutex
	t    map[uuid.UUID]Task
	adds chan Task
}

func NewPool() *Pool {
	return &Pool{
		mu:   new(sync.Mutex),
		t:    make(map[uuid.UUID]Task),
		adds: make(chan Task, 10),
	}
}

// Get returns a task of a specific type and id from the pool.
func Get[T Task](tasks *Pool, id uuid.UUID) T {
	tasks.mu.Lock()
	defer tasks.mu.Unlock()

	t, ok := tasks.t[id].(T)
	if !ok {
		var zero T
		return zero
	}
	return t
}

// Collect returns a slice of tasks of a specific type from the pool.
func Collect[T Task](tasks *Pool) []T {
	tasks.mu.Lock()
	defer tasks.mu.Unlock()

	ts := make([]T, 0, len(tasks.t))
	for _, t := range tasks.t {
		tt, ok := t.(T)
		if !ok {
			continue
		}
		ts = append(ts, tt)
	}
	return ts
}

func Contains[T Task](tasks *Pool) bool {
	tasks.mu.Lock()
	defer tasks.mu.Unlock()

	for _, t := range tasks.t {
		_, is := t.(T)
		if is {
			return true
		}
	}
	return false
}

// Get returns a task from the pool.
func (p *Pool) Get(id uuid.UUID) Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.t[id]
}

// Add adds a task to the pool.
func (p *Pool) Add(t Task) {
	slog.Debug("adding task to a pool", "type", fmt.Sprintf("%T", t), "task", t)

	p.mu.Lock()
	defer p.mu.Unlock()

	p.t[t.ID()] = t
	p.adds <- t
}

// Adds returns a channel that receives tasks that are added to the pool.
func (p *Pool) NotifyAdd() <-chan Task {
	return p.adds
}

// Remove removes a task from the pool.
func (p *Pool) Remove(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.t, t.ID())
}

// Tasks returns a copy of all tasks in the pool.
func (p *Pool) Tasks() []Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	tasks := make([]Task, 0, len(p.t))
	for _, t := range p.t {
		tasks = append(tasks, t)
	}
	return tasks
}

// Evaluate returns true if all tasks in the pool are of type Evaluate.
func (p *Pool) Evaluating() bool {
	return len(Collect[*Evaluate](p)) > 0
}

// Punctual returns true if all tasks in the pool did not exceed their deadlines.
func (p *Pool) Punctual() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for _, t := range p.t {
		if t.Deadline().Before(now) {
			return false
		}
	}
	return true
}

func (p *Pool) Idle() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, t := range p.t {
		if t.Deadline().Before(time.Now()) {
			return false
		}
		switch t.(type) {
		case *Evaluate, *Initialize, *Optimize, *ShareState:
			return false
		}
	}
	return true
}
