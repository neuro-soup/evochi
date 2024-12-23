package service

import (
	"context"
	"fmt"
	"log/slog"

	connect "github.com/bufbuild/connect-go"
	"github.com/neuro-soup/evochi/server/internal/event"
	"github.com/neuro-soup/evochi/server/internal/worker"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var subscriptions = promauto.NewCounter(prometheus.CounterOpts{
	Name: "service_subscribe_total",
	Help: "The total number of subscribe requests handled.",
})

func (h *Handler) Subscribe(
	ctx context.Context,
	req *connect.Request[evochiv1.SubscribeRequest],
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
) error {
	slog.Debug("handling subscribe request", "req", req.Msg)
	subscriptions.Inc()

	// create worker and add it to the pool
	w := worker.NewWorker(uint(req.Msg.Cores), req.Msg.Attrs)
	h.workers.Add(w)

	// send hello event to client
	h.events.Push(w, &event.Hello{ID: w.ID})

	// event loop
loop:
	for {
		select {
		case <-ctx.Done():
			// client disconnected
			slog.Debug("worker disconnected", "worker", w.ID)
			break loop

		case evt := <-h.events.Pull(w):
			// send event to client
			slog.Debug("sending event to client", "worker", w.ID, "event", evt)

			err := stream.Send(eventToProto(evt))
			if err != nil {
				slog.Error("failed to send event to client", "worker", w.ID, "event", evt, "err", err)
				break loop
			}
		}
	}

	slog.Debug("closing worker", "worker", w.ID)
	return nil
}

// eventToProto converts event to proto message
func eventToProto(evt event.Event) *evochiv1.SubscribeResponse {
	switch evt := evt.(type) {
	case *event.Hello:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_HELLO,
			Event: &evochiv1.SubscribeResponse_Hello{
				Hello: &evochiv1.HelloEvent{
					Id: evt.ID.String(),
				},
			},
		}

	default:
		panic(fmt.Sprintf("unsupported event type: %T", evt))
	}
}
