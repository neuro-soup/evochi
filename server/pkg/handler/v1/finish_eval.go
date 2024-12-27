package v1

import (
	"context"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishEvaluation(
	_ context.Context,
	req *connect.Request[evochiv1.FinishEvaluationRequest],
) (*connect.Response[evochiv1.FinishEvaluationResponse], error) {
	slog.Debug("handling finish evaluation request",
		"task", req.Msg.TaskId,
		"rewards", len(req.Msg.Rewards),
	)

	panic("unimplemented") // TODO: implement
}
