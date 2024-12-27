package event

import (
	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/epoch"
)

type Event interface {
	isEvent()
}

type Hello struct {
	// ID is the new identifier for the worker who received the event.
	ID uuid.UUID

	// PopulationSize is the total population size.
	PopulationSize uint

	// State is the initial state of the worker. If nil, the worker is the first
	// one to join the work force.
	State []byte

	// Attrs is the custom attributes set by the server.
	Attrs map[string][]byte
}

func (Hello) isEvent() {}

type Evaluate struct {
	// TaskID is the identifier of the evaluation.
	TaskID uuid.UUID

	// Slices is the slices to be evaluated.
	Slices []epoch.Slice
}

func (Evaluate) isEvent() {}

type Optimize struct {
	// TaskID is the identifier of the evaluation.
	TaskID uuid.UUID

	// Epoch is the current epoch.
	Epoch uint

	// Rewards is the accumulated rewards in the current epoch.
	Rewards [][]byte

	// ShareState is whether to send the optimized state to the server after
	// the optimization step.
	ShareState bool
}

func (Optimize) isEvent() {}
