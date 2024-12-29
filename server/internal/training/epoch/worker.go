package epoch

import "github.com/google/uuid"

type worker interface {
	// WorkerID returns the ID of this worker.
	WorkerID() uuid.UUID

	// WorkerCores returns the number of cores this worker has.
	WorkerCores() uint
}

type dummyWorker struct {
	id    uuid.UUID
	cores uint
}

func newWorker(cores uint) worker {
	return &dummyWorker{
		id:    uuid.New(),
		cores: cores,
	}
}

func (w *dummyWorker) WorkerID() uuid.UUID { return w.id }
func (w *dummyWorker) WorkerCores() uint   { return w.cores }
