package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"github.com/neuro-soup/evochi/server/internal/distribution/task"
	"github.com/neuro-soup/evochi/server/internal/training/epoch"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishInitialization(
	ctx context.Context,
	req *connect.Request[evochiv1.FinishInitializationRequest],
) (*connect.Response[evochiv1.FinishInitializationResponse], error) {
	slog.Debug("handling initialization request",
		"task", req.Msg.TaskId,
		"state", len(req.Msg.State),
	)

	// authenticate worker
	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	// check if epoch has been created
	if h.epoch != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"epoch %d has already been created", h.epoch.Number,
		))
	}

	// check if epoch has been initialized
	if h.epoch.State != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"epoch %d has already been initialized",
			h.epoch.Number,
		))
	}

	// parse task id
	taskID, err := uuid.Parse(req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"invalid task id: %w", err,
		))
	}

	// check if task exists
	t := task.Get[*task.Initialize](w.Tasks, taskID)
	if t == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf(
			"task %s not found", taskID,
		))
	}

	// check if task is related to the current epoch
	if t.Epoch != h.epoch.Number {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"task %s is not related to epoch %d", taskID, h.epoch.Number,
		))
	}

	// create new epoch
	h.epoch = epoch.New(t.Epoch, h.cfg.PopulationSize, req.Msg.State)

	// start evaluation
	for _, a := range h.workers.Workers() {
		// try to assign slices to the worker
		slices := h.epoch.Assign(a)
		if len(slices) == 0 {
			// no slices assigned
			continue
		}

		// add task to worker
		t := task.NewEvaluate(h.epoch.Number, slices, h.cfg.WorkerTimeout)
		w.Tasks.Add(t)
	}

	return connect.NewResponse(&evochiv1.FinishInitializationResponse{Ok: true}), nil
}
