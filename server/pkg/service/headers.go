package service

import (
	"net/http"

	"github.com/google/uuid"
)

// workerID returns worker ID from the request header.
func workerID(h http.Header) (uuid.UUID, error) {
	raw := h.Get("Evochi-Worker-ID")
	if raw == "" {
		return uuid.Nil, nil
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
