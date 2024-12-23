package epoch

import "time"

type Epoch struct {
	// Number is the epoch number.
	Number uint

	// StartTime is the start time of the epoch.
	StartTime time.Time
}

// New creates a new epoch.
func New(number uint, population uint) *Epoch {
	return &Epoch{
		Number:    number,
		StartTime: time.Now(),
	}
}
