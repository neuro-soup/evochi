package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/distribution/task"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishShareState(
	_ context.Context,
	req *connect.Request[evochiv1.FinishShareStateRequest],
) (*connect.Response[evochiv1.FinishShareStateResponse], error) {
	slog.Debug("handling finish share state request",
		"task", req.Msg.TaskId,
		"state", len(req.Msg.State),
	)

	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	taskID, err := uuid.Parse(req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"invalid task id %w", err,
		))
	}

	t := task.Get[*task.ShareState](w.Tasks, taskID)
	if t == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf(
			"task %s not found", req.Msg.TaskId,
		))
	}

	// complete task
	t.Done()
	w.Tasks.Remove(t)

	if t.Epoch == h.epoch.Number && h.finished() {
		// all workers have finished their state sharing
		// go to next epoch
		h.nextEpoch(req.Msg.State)
	}

	msg := &evochiv1.FinishShareStateResponse{Ok: true}
	return &connect.Response[evochiv1.FinishShareStateResponse]{Msg: msg}, nil
}
