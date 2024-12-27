package worker

import "time"

type heartbeat struct {
	SeqID      uint      `json:"seq_id"`
	SentAt     time.Time `json:"sent_at"`
	ReceivedAt time.Time `json:"received_at"`
}

func (h heartbeat) Ping() time.Duration {
	return h.ReceivedAt.Sub(h.SentAt)
}
