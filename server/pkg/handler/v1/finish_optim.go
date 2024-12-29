package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishOptimization(
	_ context.Context,
	req *connect.Request[evochiv1.FinishOptimizationRequest],
) (*connect.Response[evochiv1.FinishOptimizationResponse], error) {
	slog.Debug("handling finish optimization request",
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

	// parse task id
	taskID, err := uuid.Parse(req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf(
			"invalid task id: %w", err,
		))
	}

	_ = taskID
	_ = w

	panic("unimplemented") // TODO: implement
}
