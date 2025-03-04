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

func (h *Handler) FinishOptimization(
	_ context.Context,
	req *connect.Request[evochiv1.FinishOptimizationRequest],
) (*connect.Response[evochiv1.FinishOptimizationResponse], error) {
	slog.Debug("handling finish optimization request",
		"task", req.Msg.TaskId,
	)

	// authenticate worker
	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	// check if epoch is initialised
	if h.epoch == nil || h.epoch.State == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"epoch is not initialised",
		))
	}

	// parse task id
	taskID, err := uuid.Parse(req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"invalid task id: %w", err,
		))
	}

	t := task.Get[*task.Optimize](w.Tasks, taskID)
	if t == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf(
			"task %s not found", taskID,
		))
	}

	// complete task
	t.Done()
	w.Tasks.Remove(t)

	if h.finished() {
		// all workers have finished their optimization
		// request optimized state from trusted workers
		h.requestState()
	}

	return connect.NewResponse(&evochiv1.FinishOptimizationResponse{Ok: true}), nil
}
