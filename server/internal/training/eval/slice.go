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

func (s Slice) Overlaps(other Slice) bool {
	return s.Start < other.End && other.Start < s.End
}

func SliceWidth(slices []Slice) uint {
	var total uint
	for _, s := range slices {
		total += s.Width()
	}
	return total
}

func SliceOverlaps(slices ...Slice) *[2]Slice {
	for i, s := range slices {
		for j, other := range slices {
			if i != j && s.Overlaps(other) {
				return &[2]Slice{s, other}
			}
		}
	}
	return nil
}
