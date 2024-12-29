package epoch

import (
	"fmt"
	"log/slog"

	"github.com/neuro-soup/evochi/server/internal/stack"
	"github.com/neuro-soup/evochi/server/internal/training/eval"
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
	State eval.State

	// unassigned is the stack of slices that are not assigned to any worker.
	unassigned *stack.Stack[eval.Slice]

	// rewards are the collected rewards for this epoch. Its size is equal to
	// the population.
	rewards []eval.Reward
}

// New creates a new epoch.
func New(number, population uint, initState eval.State) *Epoch {
	epochsCreated.Inc()
	return &Epoch{
		Number:     number,
		Population: population,
		State:      initState,
		unassigned: stack.New(
			eval.Slice{
				Start: 0,
				End:   population,
			},
		),
		rewards: make([]eval.Reward, population),
	}
}

func (e *Epoch) Assign(w worker) []eval.Slice {
	workerID, workerCores := w.WorkerID(), w.WorkerCores()

	if e.unassigned.Len() == 0 {
		slog.Debug("no slices to assign",
			"epoch", e.Number,
			"worker", workerID,
			"cores", workerCores,
		)
		return nil
	}

	var (
		slices []eval.Slice
		width  uint
	)

	for width < workerCores && e.unassigned.Len() > 0 {
		pop := e.unassigned.Pop()
		delta := pop.End - pop.Start

		if delta <= workerCores-width {
			// assign the whole slice to the worker
			slices = append(slices, pop)
			width += delta
			continue
		}

		// assign the first part of the slice to the worker
		first := eval.Slice{
			Start: pop.Start,
			End:   pop.Start + min(workerCores-width, delta),
		}
		slices = append(slices, first)
		width += first.End - first.Start

		// push the rest of the slice to the stack
		second := eval.Slice{
			Start: first.End,
			End:   pop.End,
		}
		if second.Start != second.End {
			e.unassigned.Push(second)
		}
	}

	return slices
}

func (e *Epoch) Reward(w worker, slices []eval.Slice, rewards []eval.Reward) error {
	workerID, workerCores := w.WorkerID(), w.WorkerCores()

	width := eval.TotalSliceWidth(slices)
	if len(rewards) != int(width) {
		slog.Error("invalid number of rewards",
			"epoch", e.Number,
			"worker", workerID,
			"slices", width,
			"got", len(rewards),
		)
		return fmt.Errorf("expected %d rewards, got %d", width, len(rewards))
	}

	slog.Debug("rewarding epoch for slices",
		"worker", workerID,
		"cores", workerCores,
		"slices", len(slices),
		"rewards", len(rewards),
	)

	for _, slice := range slices {
		for i := slice.Start; i < slice.End; i++ {
			e.rewards[i] = rewards[i-slice.Start]
		}
	}

	return nil
}
