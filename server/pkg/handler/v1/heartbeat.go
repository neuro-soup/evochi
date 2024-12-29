package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/neuro-soup/evochi/server/internal/distribution/task"
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

	// authenticate worker
	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	// check if there is a heartbeat task
	hbs := task.Collect[*task.Heartbeat](w.Tasks)
	if len(hbs) == 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"no heartbeat tasks found",
		))
	}

	hb := hbs[0]
	if hb.SeqID != uint(req.Msg.SeqId) {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"heartbeat mismatch: expected %d, got %d",
			hb.SeqID, req.Msg.SeqId,
		))
	}

	// complete the task
	hb.Done()
	w.Tasks.Remove(hb)

	// add new heartbeat task
	w.Tasks.Add(task.NewHeartbeat(hb.SeqID+1, h.cfg.WorkerTimeout))

	return connect.NewResponse(&evochiv1.HeartbeatResponse{Ok: true}), nil
}
