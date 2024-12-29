package eval

type Slice struct {
	// Start is the start index (inclusive).
	Start uint

	// End is the end index (exclusive).
	End uint
}

func (s Slice) Width() uint {
	return s.End - s.Start
}

func TotalSliceWidth(slices []Slice) uint {
	var total uint
	for _, s := range slices {
		total += s.Width()
	}
	return total
}
