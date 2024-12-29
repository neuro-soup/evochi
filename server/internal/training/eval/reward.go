package eval

type Reward []byte

func BytesToRewards(b [][]byte) []Reward {
	r := make([]Reward, len(b))
	for i, b := range b {
		r[i] = Reward(b)
	}
	return r
}
