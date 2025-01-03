package epoch

import (
	"testing"

	"github.com/neuro-soup/evochi/server/internal/training/eval"
	"github.com/stretchr/testify/require"
)

func TestEpoch_Assign(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 100, nil)
		w := newWorker(4)

		e.unassigned.Clear()

		assigned := e.Assign(w)

		r.Empty(assigned)
	})

	t.Run("Exact", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := newWorker(10)

		assigned := e.Assign(w)

		r.Len(assigned, 1)
		r.EqualValues(0, assigned[0].Start)
		r.EqualValues(10, assigned[0].End)

		r.EqualValues(0, e.unassigned.Len())
	})

	t.Run("Exhaustive Single Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 5, nil)
		w := newWorker(5)

		assigned := e.Assign(w)

		r.Len(assigned, 1)
		r.EqualValues(0, assigned[0].Start)
		r.EqualValues(5, assigned[0].End)

		r.EqualValues(0, e.unassigned.Len())
	})

	t.Run("Inexhaustive Single Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := newWorker(5)

		assigned := e.Assign(w)

		r.Len(assigned, 1)
		r.EqualValues(0, assigned[0].Start)
		r.EqualValues(5, assigned[0].End)

		r.EqualValues(1, e.unassigned.Len())
		r.EqualValues(5, e.unassigned.Peek().Start)
		r.EqualValues(10, e.unassigned.Peek().End)
	})

	t.Run("Exhaustive Multiple Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := newWorker(10)

		e.unassigned.Clear()
		e.unassigned.Push(eval.Slice{
			Start: 0,
			End:   5,
		})
		e.unassigned.Push(eval.Slice{
			Start: 7,
			End:   10,
		})

		assigned := e.Assign(w)

		r.Len(assigned, 2)
		r.EqualValues(7, assigned[0].Start)
		r.EqualValues(10, assigned[0].End)
		r.EqualValues(0, assigned[1].Start)
		r.EqualValues(5, assigned[1].End)

		r.EqualValues(0, e.unassigned.Len())
	})

	t.Run("Inexhaustive Multiple Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 7, nil)
		w := newWorker(5)

		e.unassigned.Clear()
		e.unassigned.Push(eval.Slice{
			Start: 0,
			End:   5,
		})
		e.unassigned.Push(eval.Slice{
			Start: 7,
			End:   10,
		})

		assigned := e.Assign(w)

		r.Len(assigned, 2)
		r.EqualValues(7, assigned[0].Start)
		r.EqualValues(10, assigned[0].End)
		r.EqualValues(0, assigned[1].Start)
		r.EqualValues(2, assigned[1].End)

		r.EqualValues(1, e.unassigned.Len())
		r.EqualValues(2, e.unassigned.Peek().Start)
		r.EqualValues(5, e.unassigned.Peek().End)
	})
}

func TestEpoch_Reward(t *testing.T) {
	t.Run("Overlapping", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := newWorker(10)

		e.unassigned.Clear()

		evals := []eval.Eval{
			{
				Slice: eval.Slice{
					Start: 0,
					End:   2,
				},
				Rewards: []float64{1, 1},
			},
			{
				Slice: eval.Slice{
					Start: 0,
					End:   3,
				},
				Rewards: []float64{1, 1, 1},
			},
		}
		err := e.Reward(w, evals)

		r.Error(err)
	})

	t.Run("Reward", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := newWorker(10)

		e.unassigned.Clear()

		evals := []eval.Eval{
			{
				Slice: eval.Slice{
					Start: 0,
					End:   2,
				},
				Rewards: []float64{1, 2},
			},
			{
				Slice: eval.Slice{
					Start: 3,
					End:   5,
				},
				Rewards: []float64{3, 4},
			},
		}
		err := e.Reward(w, evals)

		r.NoError(err)
		r.EqualValues(1, e.rewards[0])
		r.EqualValues(2, e.rewards[1])
		r.EqualValues(3, e.rewards[3])
		r.EqualValues(4, e.rewards[4])
		for i, rew := range e.rewards {
			if i == 0 || i == 1 || i == 3 || i == 4 {
				continue
			}
			r.EqualValues(0, rew, "expected nor reward at %d", i)
		}
	})
}
