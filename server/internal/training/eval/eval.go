package eval

type Eval struct {
	Slice   Slice
	Rewards []float64
}

func EvalSlices(evals ...Eval) []Slice {
	var slices []Slice
	for _, e := range evals {
		slices = append(slices, e.Slice)
	}
	return slices
}
