package task

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	mu    *sync.Mutex
	tasks map[uuid.UUID]Task
}

func NewPool() *Pool {
	return &Pool{
		mu:    new(sync.Mutex),
		tasks: make(map[uuid.UUID]Task),
	}
}

func (p *Pool) Get(id uuid.UUID) Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.tasks[id]
}

func (p *Pool) Add(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.tasks[t.ID()] = t
}

func (p *Pool) Remove(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.tasks, t.ID())
}

func (p *Pool) Tasks() []Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	tasks := make([]Task, 0, len(p.tasks))
	for _, t := range p.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (p *Pool) Punctual() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for _, t := range p.tasks {
		if t.Deadline().After(now) {
			return false
		}
	}
	return true
}
