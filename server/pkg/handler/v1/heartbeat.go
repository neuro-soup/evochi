package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) Heartbeat(
	_ context.Context,
	req *connect.Request[evochiv1.HeartbeatRequest],
) (*connect.Response[evochiv1.HeartbeatResponse], error) {
	slog.Debug("handling heartbeat request",
		"seq_id", req.Msg.SeqId,
		"timestamp", req.Msg.Timestamp,
	)

	worker, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	err = worker.Heartbeat(uint(req.Msg.SeqId), req.Msg.Timestamp.AsTime())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf(
			"heartbeat failed: %w", err,
		))
	}

	return connect.NewResponse(&evochiv1.HeartbeatResponse{Ok: true}), nil
}
