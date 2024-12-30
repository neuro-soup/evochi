package eval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice_Overlaps(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		r := require.New(t)

		s1 := Slice{Start: 0, End: 0}
		s2 := Slice{Start: 0, End: 0}

		r.False(s1.Overlaps(s2))
	})

	t.Run("Overlap", func(t *testing.T) {
		r := require.New(t)

		s1 := Slice{Start: 0, End: 10}
		s2 := Slice{Start: 5, End: 15}

		r.True(s1.Overlaps(s2))
	})

	t.Run("Disjoint", func(t *testing.T) {
		r := require.New(t)

		s1 := Slice{Start: 0, End: 10}
		s2 := Slice{Start: 10, End: 20}

		r.False(s1.Overlaps(s2))
	})
}

func TestSliceOverlaps(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		r := require.New(t)

		s := []Slice{
			{
				Start: 0,
				End:   2,
			},
		}

		overlap := SliceOverlaps(s...)
		r.Nil(overlap)
	})
}
