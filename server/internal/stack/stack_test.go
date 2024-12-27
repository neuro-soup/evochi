package stack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		r := require.New(t)

		st := New[int]()

		r.Panics(func() { _ = st.Pop() })
	})

	t.Run("Single", func(t *testing.T) {
		r := require.New(t)

		st := New[int]()

		st.Push(1)

		r.Equal(1, st.Pop())
		r.Panics(func() { _ = st.Pop() })
	})

	t.Run("Two", func(t *testing.T) {
		r := require.New(t)

		st := New[int]()

		st.Push(1)
		st.Push(2)

		r.Equal(2, st.Pop())
		r.Equal(1, st.Pop())
		r.Panics(func() { _ = st.Pop() })
	})
}
