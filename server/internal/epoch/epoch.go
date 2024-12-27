package epoch

import (
	"log/slog"

	"github.com/neuro-soup/evochi/server/internal/stack"
	"github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var epochsCreated = promauto.NewCounter(prometheus.CounterOpts{
	Name: "epochs_created",
	Help: "The number of epochs created.",
})

type Epoch struct {
	// Number is the epoch number.
	Number uint

	// Population is the size of the population to be evaluated.
	Population uint

	// State is the state used when starting the epoch.
	State []byte

	// unassigned is the stack of slices that are not assigned to any worker.
	unassigned *stack.Stack[Slice]
}

// New creates a new epoch.
func New(number, population uint, initState []byte) *Epoch {
	epochsCreated.Inc()
	return &Epoch{
		Number:     number,
		Population: population,
		State:      initState,
		unassigned: stack.New(
			Slice{
				Start: 0,
				End:   population,
			},
		),
	}
}

func (e *Epoch) Assign(w *worker.Worker) []Slice {
	if e.unassigned.Len() == 0 {
		slog.Debug("no slices to assign", "epoch", e.Number, "worker", w.ID)
		return nil
	}

	var (
		slices []Slice
		width  uint
	)

	for width < w.Cores && e.unassigned.Len() > 0 {
		pop := e.unassigned.Pop()
		delta := pop.End - pop.Start

		if delta <= w.Cores-width {
			// assign the whole slice to the worker
			slices = append(slices, pop)
			width += delta
			continue
		}

		// assign the first part of the slice to the worker
		first := Slice{
			Start: pop.Start,
			End:   pop.Start + min(w.Cores-width, delta),
		}
		slices = append(slices, first)
		width += first.End - first.Start

		// push the rest of the slice to the stack
		second := Slice{
			Start: first.End,
			End:   pop.End,
		}
		if second.Start != second.End {
			e.unassigned.Push(second)
		}
	}

	return slices
}
