package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	connect "github.com/bufbuild/connect-go"
	"github.com/neuro-soup/evochi/server/internal/distribution/task"
	"github.com/neuro-soup/evochi/server/internal/distribution/worker"
	"github.com/neuro-soup/evochi/server/internal/training/epoch"
	"github.com/neuro-soup/evochi/server/internal/training/eval"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func (h *Handler) Subscribe(
	ctx context.Context,
	req *connect.Request[evochiv1.SubscribeRequest],
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
) error {
	slog.Debug("handling subscribe request",
		"cores", req.Msg.Cores,
	)

	// create worker
	w := worker.New(uint(req.Msg.Cores))

	// add worker to pools and cleanup on exit
	h.workers.Add(w)
	defer h.workers.Remove(w)

	// demand heartbeat
	w.Tasks.Add(task.NewHeartbeat(1, h.cfg.WorkerTimeout))

	if h.epoch == nil {
		// first worker, create and initialise epoch, send hello event
		h.handleFirstSubscriber(stream, w)
	} else {
		// subsequent worker, send hello event
		h.handleSubsequentSubscriber(stream, w)
	}

	// event-loop
	for {
		select {
		case <-w.Removes(): // server-side cancellation
			slog.Debug("server cancelled subscription")

			h.handleCancellation(w)
			return nil

		case <-ctx.Done(): // client-side cancellation
			slog.Debug("client cancelled subscription")

			h.handleCancellation(w)
			return ctx.Err()

		case t := <-w.Tasks.Notify():
			resp := taskToProto(t)
			if resp == nil {
				// ignore tasks that are not of interest
				continue
			}

			slog.Debug("sending task to worker", "event", t, "worker", w.ID)

			if err := stream.Send(resp); err != nil {
				slog.Error("failed to send task to worker", "error", err)

				h.handleCancellation(w)
				return err
			}
		}
	}
}

func (h *Handler) handleFirstSubscriber(
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
	w *worker.Worker,
) {
	slog.Debug("handling first subscriber", "worker", w.ID)

	h.epoch = epoch.New(1, h.cfg.PopulationSize, nil)
	t := task.NewInitialize(h.epoch.Number, h.cfg.WorkerTimeout)

	// initialise worker by sending hello event and demand initial state
	err := stream.Send(&evochiv1.SubscribeResponse{
		Event: &evochiv1.SubscribeResponse_Hello{
			Hello: &evochiv1.HelloEvent{
				Id:             w.ID.String(),
				PopulationSize: int32(h.cfg.PopulationSize),
				State:          nil, // is initialised and set by the client
			},
		},
	})
	if err != nil {
		slog.Error("failed to send hello event", "error", err)
		w.Remove()
		return
	}

	w.Tasks.Add(t)
}

func (h *Handler) handleSubsequentSubscriber(
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
	w *worker.Worker,
) {
	slog.Debug("handling subsequent subscriber", "worker", w.ID)

	// initialise worker by sending hello event
	err := stream.Send(&evochiv1.SubscribeResponse{
		Event: &evochiv1.SubscribeResponse_Hello{
			Hello: &evochiv1.HelloEvent{
				Id:             w.ID.String(),
				PopulationSize: int32(h.epoch.Population),
				State:          h.epoch.State,
			},
		},
	})
	if err != nil {
		slog.Error("failed to send hello event", "error", err)
		w.Remove()
		return
	}

	if h.epoch == nil || h.epoch.State == nil {
		// epoch is not initialised yet
		return
	}

	// try to assign slices to the worker
	assigned := h.epoch.Assign(w)
	if len(assigned) == 0 {
		// no tasks to assign, worker is idle
		return
	}

	// add task to worker, which will be picked up by the event-loop
	t := task.NewEvaluate(h.epoch.Number, assigned, h.cfg.WorkerTimeout)
	w.Tasks.Add(t)
}

func (h *Handler) handleCancellation(w *worker.Worker) {
	suc := h.workers.Trusted(func(f *worker.Worker) bool { return f.ID != w.ID })
	if suc == nil {
		// no successor worker available
		h.epoch = nil
		return
	}

	// re-distribute work load of cancelled worker
	for _, t := range w.Tasks.Tasks() {
		switch t := t.(type) {
		case *task.Initialize:
			t.RequestedAt = time.Now()
			suc.Tasks.Add(t)
		case *task.Evaluate:
			t.RequestedAt = time.Now()
			suc.Tasks.Add(t)
		}
	}

	// remove worker from pool
	w.Remove()
}

func taskToProto(t task.Task) *evochiv1.SubscribeResponse {
	switch t := t.(type) {
	case *task.Initialize:
		return &evochiv1.SubscribeResponse{
			Event: &evochiv1.SubscribeResponse_Initialize{
				Initialize: &evochiv1.InitializeEvent{
					TaskId: t.ID().String(),
				},
			},
		}

	case *task.Evaluate:
		return &evochiv1.SubscribeResponse{
			Event: &evochiv1.SubscribeResponse_Evaluate{
				Evaluate: &evochiv1.EvaluateEvent{
					TaskId: t.ID().String(),
					Slices: slicesToProto(t.Slices),
				},
			},
		}

	case *task.Heartbeat:
		// ignore heartbeats
		return nil

	default:
		panic(fmt.Sprintf("unknown task type: %T", t))
	}
}

func slicesToProto(slices []eval.Slice) []*evochiv1.Slice {
	out := make([]*evochiv1.Slice, len(slices))
	for i, slice := range slices {
		out[i] = &evochiv1.Slice{
			Start: int32(slice.Start),
			End:   int32(slice.End),
		}
	}
	return out
}
