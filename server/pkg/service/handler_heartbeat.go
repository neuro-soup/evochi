package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	connect "github.com/bufbuild/connect-go"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var heartbeats = promauto.NewCounter(prometheus.CounterOpts{
	Name: "service_heartbeat_total",
	Help: "The total number of heartbeat requests handled.",
})

func (h *Handler) Heartbeat(
	ctx context.Context,
	req *connect.Request[evochiv1.HeartbeatRequest],
) (*connect.Response[evochiv1.HeartbeatResponse], error) {
	slog.Debug("handling heartbeat request", "req", req.Msg)
	heartbeats.Inc()

	workerID, err := workerID(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"invalid worker id: %w", err,
		))
	}

	w := h.workers.Get(workerID)
	if w == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf(
			"worker %q not found", workerID,
		))
	}

	w.Ping = time.Since(req.Msg.Timestamp.AsTime())
	w.LastSeen = time.Now()

	resp := &evochiv1.HeartbeatResponse{Ok: true}
	return &connect.Response[evochiv1.HeartbeatResponse]{Msg: resp}, nil
}
