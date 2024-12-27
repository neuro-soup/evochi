package v1

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/neuro-soup/evochi/server/internal/epoch"
	"github.com/neuro-soup/evochi/server/internal/event"
	"github.com/neuro-soup/evochi/server/internal/worker"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) Subscribe(
	ctx context.Context,
	req *connect.Request[evochiv1.SubscribeRequest],
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
) error {
	slog.Debug("handling subscribe request",
		"cores", req.Msg.Cores,
		"attrs", req.Msg.Attrs,
	)

	w := worker.New(uint(req.Msg.Cores), req.Msg.Attrs)

	// add worker to pools
	h.workers.Add(w)
	h.conns.add(w)

	// cleanup on exit
	defer h.workers.Remove(w)
	defer h.events.Close(w)
	defer h.conns.remove(w)

	// TODO: assign evals if epoch is not nil

	// event-loop
	for {
		select {
		case <-h.conns.disconnects(w): // server-side cancellation
			slog.Debug("server cancelled subscription")
			// TODO: re-distribute work load of cancelled worker

			return nil

		case <-ctx.Done(): // client-side cancellation
			slog.Debug("client cancelled subscription")
			// TODO: re-distribute work load of cancelled worker

			return ctx.Err()

		case evt := <-h.events.Pull(w):
			slog.Debug("sending event to worker", "event", evt, "worker", w.ID)

			if err := stream.Send(eventToProto(evt)); err != nil {
				slog.Error("failed to send event to worker", "error", err)

				// TODO: re-distribute work load of cancelled worker
				return err
			}
		}
	}
}

// eventToProto converts an event to a proto message.
func eventToProto(evt event.Event) *evochiv1.SubscribeResponse {
	switch evt := evt.(type) {
	case event.Hello:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_HELLO,
			Event: &evochiv1.SubscribeResponse_Hello{
				Hello: &evochiv1.HelloEvent{
					Id:             evt.ID.String(),
					PopulationSize: int32(evt.PopulationSize),
					State:          evt.State,
					Attrs:          evt.Attrs,
				},
			},
		}

	case event.Evaluate:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_EVALUATE,
			Event: &evochiv1.SubscribeResponse_Evaluate{
				Evaluate: &evochiv1.EvaluateEvent{
					EvalId: evt.EvalID.String(),
					Slices: slicesToProto(evt.Slices),
				},
			},
		}

	case event.Optimize:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_OPTIMIZE,
			Event: &evochiv1.SubscribeResponse_Optimize{
				Optimize: &evochiv1.OptimizeEvent{
					Epoch:      int32(evt.Epoch),
					Rewards:    evt.Rewards,
					ShareState: evt.ShareState,
				},
			},
		}

	default:
		panic(fmt.Errorf("unknown event type: %T", evt))
	}
}

func slicesToProto(slices []epoch.Slice) []*evochiv1.Slice {
	out := make([]*evochiv1.Slice, len(slices))
	for i, slice := range slices {
		out[i] = &evochiv1.Slice{
			Start: int32(slice.Start),
			End:   int32(slice.End),
		}
	}
	return out
}
