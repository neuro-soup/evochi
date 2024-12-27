package task

import (
	"time"

	"github.com/google/uuid"
)

type Task interface {
	// ID returns the unique identifier of the task.
	ID() uuid.UUID

	// Deadline returns the time by which the task must be completed. If the worker
	// does not complete the task before this time, the task is considered failed
	// and the worker removed.
	Deadline() time.Time

	// Done is called when the task is completed.
	Done()
}
