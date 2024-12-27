package v1

import (
	"sync"

	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/worker"
)

type connections struct {
	mu              *sync.Mutex
	disconnectChans map[uuid.UUID]chan struct{}
}

func newConnections() *connections {
	return &connections{
		mu:              new(sync.Mutex),
		disconnectChans: make(map[uuid.UUID]chan struct{}),
	}
}

func (c *connections) disconnect(w *worker.Worker) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := c.disconnectChans[w.ID]
	if ch == nil {
		return
	}
	ch <- struct{}{}
	delete(c.disconnectChans, w.ID)
}

func (c *connections) disconnects(w *worker.Worker) <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := c.disconnectChans[w.ID]
	if ch == nil {
		ch = make(chan struct{})
		c.disconnectChans[w.ID] = ch
	}
	return ch
}

func (c *connections) add(w *worker.Worker) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.disconnectChans[w.ID] = make(chan struct{})
}

func (c *connections) remove(w *worker.Worker) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.disconnectChans, w.ID)
}
