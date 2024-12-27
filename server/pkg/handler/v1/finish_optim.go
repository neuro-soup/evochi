package v1

import (
	"context"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) FinishOptimization(
	_ context.Context,
	req *connect.Request[evochiv1.FinishOptimizationRequest],
) (*connect.Response[evochiv1.FinishOptimizationResponse], error) {
	slog.Debug("handling finish optimization request",
		"state_size", len(req.Msg.State),
	)

	panic("unimplemented") // TODO: implement
}
