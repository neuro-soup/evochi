package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/neuro-soup/evochi/server/internal/worker/task"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishInitialization(
	ctx context.Context,
	req *connect.Request[evochiv1.FinishInitializationRequest],
) (*connect.Response[evochiv1.FinishInitializationResponse], error) {
	slog.Debug("handling initialization request",
		"task", req.Msg.TaskId,
		"state_size", len(req.Msg.State),
	)

	w, _, err := h.authenticateWorker(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(
			"failed to authenticate worker: %w", err,
		))
	}

	if h.epoch == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"no epoch has been created yet",
		))
	}

	if h.epoch.State != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"epoch %d has already been initialized",
			h.epoch.Number,
		))
	}

	inits := task.Initializes(w.Tasks)
	if len(inits) == 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"no initialization tasks have been created yet",
		))
	}
	init := inits[0]
	if init.Epoch != h.epoch.Number {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf(
			"initialization task for epoch %d has been created, but epoch %d is being initialized",
			init.Epoch, h.epoch.Number,
		))
	}

	h.epoch.State = req.Msg.State

	return connect.NewResponse(&evochiv1.FinishInitializationResponse{Ok: true}), nil
}
