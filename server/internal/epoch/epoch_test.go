package epoch

import (
	"testing"

	"github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/stretchr/testify/require"
)

func TestEpoch_Assign(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 100, nil)
		w := worker.New(4, nil)

		e.unassigned.Clear()

		assigned := e.Assign(w)

		r.Empty(assigned)
	})

	t.Run("Exact", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := worker.New(10, nil)

		assigned := e.Assign(w)

		r.Len(assigned, 1)
		r.EqualValues(0, assigned[0].Start)
		r.EqualValues(10, assigned[0].End)

		r.EqualValues(0, e.unassigned.Len())
	})

	t.Run("Exhaustive Single Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 5, nil)
		w := worker.New(5, nil)

		assigned := e.Assign(w)

		r.Len(assigned, 1)
		r.EqualValues(0, assigned[0].Start)
		r.EqualValues(5, assigned[0].End)

		r.EqualValues(0, e.unassigned.Len())
	})

	t.Run("Inexhaustive Single Partial", func(t *testing.T) {
		r := require.New(t)

		e := New(1, 10, nil)
		w := worker.New(5, nil)

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
		w := worker.New(10, nil)

		e.unassigned.Clear()
		e.unassigned.Push(Slice{
			Start: 0,
			End:   5,
		})
		e.unassigned.Push(Slice{
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
		w := worker.New(5, nil)

		e.unassigned.Clear()
		e.unassigned.Push(Slice{
			Start: 0,
			End:   5,
		})
		e.unassigned.Push(Slice{
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
