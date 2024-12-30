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

func (h *Handler) FinishEvaluation(
	_ context.Context,
	req *connect.Request[evochiv1.FinishEvaluationRequest],
) (*connect.Response[evochiv1.FinishEvaluationResponse], error) {
	slog.Debug("handling finish evaluation request",
		"task", req.Msg.TaskId,
		"evals", len(req.Msg.Evaluations),
	)

	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	// check whether epoch has been created
	if h.epoch == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"no epoch has been created yet",
		))
	}

	// check whether epoch has been initialised by a worker
	if h.epoch.State == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"epoch %d has not been initialized",
			h.epoch.Number,
		))
	}

	// parse task id
	taskID, err := uuid.Parse(req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"invalid task id %w", err,
		))
	}

	// check whether task exists
	t := task.Get[*task.Evaluate](w.Tasks, taskID)
	if t == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf(
			"task %s not found", req.Msg.TaskId,
		))
	}

	// complete task
	t.Done()
	w.Tasks.Remove(t)

	// apply rewards
	err = h.epoch.Reward(w, protoToEvals(req.Msg.Evaluations))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf(
			"failed to apply rewards: %w", err,
		))
	}

	// give next eval task
	h.eval(w)

	resp := &evochiv1.FinishEvaluationResponse{Ok: true}
	return &connect.Response[evochiv1.FinishEvaluationResponse]{Msg: resp}, nil
}
