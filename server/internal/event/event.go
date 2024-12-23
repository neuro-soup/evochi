package event

import "github.com/google/uuid"

type Event interface {
	isEvent()
}

type Hello struct {
	// ID is the new identifier for the worker who received the event.
	ID uuid.UUID
}

func (Hello) isEvent() {}
