package v1

import (
	"context"
	"errors"
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

	if int(h.cfg.MaxWorkers) > 0 && h.workers.Len() >= int(h.cfg.MaxWorkers) {
		return connect.NewError(connect.CodeResourceExhausted, fmt.Errorf(
			"max workers reached (%d)", h.cfg.MaxWorkers,
		))
	}

	// create worker
	w := worker.New(uint(req.Msg.Cores))

	// add worker to pools and cleanup on exit
	h.workers.Add(w)
	defer h.workers.Remove(w)

	// demand heartbeat
	w.Tasks.Add(task.NewHeartbeat(1, h.cfg.WorkerTimeout))

	// greet worker
	if err := h.greet(stream, w); err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	// event-loop
	return h.handleEvents(ctx, stream, w)
}

func (h *Handler) populationSize() uint {
	if h.epoch == nil {
		return h.cfg.PopulationSize
	}
	return h.epoch.Population
}

// greet initialises the worker by sending a hello event. If the worker is the
// first worker to connect, it will also create a new epoch without any state.
//
// The state must be sent by a worker by a share-state task since the server
// does not know about the structure of the state.
func (h *Handler) greet(
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
	w *worker.Worker,
) error {
	slog.Debug("greeting worker", "worker", w.ID)

	var state eval.State
	if h.epoch != nil {
		state = h.epoch.State
	}

	tok, err := h.createJWTString(w)
	if err != nil {
		slog.Error("failed to create JWT token",
			"worker", w.ID,
			"error", err,
		)
		return fmt.Errorf("failed to create JWT token: %w", err)
	}

	// initialise worker by sending hello event
	err = stream.Send(&evochiv1.SubscribeResponse{
		Type: evochiv1.EventType_EVENT_TYPE_HELLO,
		Event: &evochiv1.SubscribeResponse_Hello{
			Hello: &evochiv1.HelloEvent{
				Id:                w.ID.String(),
				Token:             tok,
				PopulationSize:    int32(h.populationSize()),
				HeartbeatInterval: int32((h.cfg.WorkerTimeout / 2).Seconds()),
				State:             state,
			},
		},
	})
	if err != nil {
		slog.Error("failed to send hello event", "error", err)
		return fmt.Errorf("failed to send hello event: %w", err)
	}

	if h.epoch == nil {
		// first worker, create and initialise epoch
		slog.Debug("first worker, creating epoch")

		h.epoch = epoch.New(1, h.cfg.PopulationSize, nil)
		t := task.NewInitialize(h.epoch.Number, h.cfg.WorkerTimeout)
		w.Tasks.Add(t)
	} else if h.epoch.State != nil {
		// subsequent worker with initialised epoch
		slog.Debug("subsequent worker, epoch started, giving work")
		h.eval(w)
	}

	return nil
}

// handleEvents is the event-loop of a worker. It manages cancellation and task
// assignment.
func (h *Handler) handleEvents(
	ctx context.Context,
	stream *connect.ServerStream[evochiv1.SubscribeResponse],
	w *worker.Worker,
) error {
	for {
		select {
		case <-w.NotifyRemoval(): // server-side cancellation
			slog.Debug("server cancelled subscription", "worker", w.ID)

			h.handleCancellation(w)
			return connect.NewError(connect.CodeCanceled, errors.New("server removed you"))

		case <-ctx.Done(): // client-side cancellation
			slog.Debug("client cancelled subscription", "worker", w.ID)

			h.handleCancellation(w)
			return ctx.Err()

		case t := <-w.Tasks.NotifyAdd(): // worker receives task
			resp := taskToProto(t)
			if resp == nil {
				// ignore tasks that are not of interest
				continue
			}

			slog.Debug("sending task to worker",
				"type", fmt.Sprintf("%T", t),
				"task", t,
				"worker", w.ID,
			)

			if err := stream.Send(resp); err != nil {
				slog.Error("failed to send task to worker", "error", err)

				h.handleCancellation(w)
				return err
			}
		}
	}
}

func (h *Handler) handleCancellation(w *worker.Worker) {
	// elect a successor worker
	// TODO: maybe prefer an idling worker
	suc := h.workers.Trusted(func(f *worker.Worker) bool { return f.ID != w.ID })
	if suc == nil {
		// no successor worker available
		slog.Debug("no successor worker available, resetting handler", "worker", w.ID)
		h.reset()
		return
	}

	// re-distribute work load of cancelled worker
	for _, t := range w.Tasks.Tasks() {
		switch t := t.(type) {
		case *task.Initialize:
			t.RequestedAt = time.Now()
			suc.Tasks.Add(t)
		case *task.Evaluate:
			h.epoch.Unassign(t.Slices...)
			// // TODO: think about whether this should be put back onto epoch
			// t.RequestedAt = time.Now()
			// suc.Tasks.Add(t)
		case *task.ShareState:
			t.RequestedAt = time.Now()
			suc.Tasks.Add(t)
		}
	}

	// remove worker from pool
	w.Remove()
}

// taskToProto converts a task to a proto subscribe response.
func taskToProto(t task.Task) *evochiv1.SubscribeResponse {
	switch t := t.(type) {
	case *task.Initialize:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_INITIALIZE,
			Event: &evochiv1.SubscribeResponse_Initialize{
				Initialize: &evochiv1.InitializeEvent{
					TaskId: t.ID().String(),
				},
			},
		}

	case *task.Evaluate:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_EVALUATE,
			Event: &evochiv1.SubscribeResponse_Evaluate{
				Evaluate: &evochiv1.EvaluateEvent{
					TaskId: t.ID().String(),
					Epoch:  int32(t.Epoch),
					Slices: slicesToProto(t.Slices),
				},
			},
		}

	case *task.ShareState:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_SHARE_STATE,
			Event: &evochiv1.SubscribeResponse_ShareState{
				ShareState: &evochiv1.ShareStateEvent{
					TaskId: t.ID().String(),
					Epoch:  int32(t.Epoch),
				},
			},
		}

	case *task.Optimize:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_OPTIMIZE,
			Event: &evochiv1.SubscribeResponse_Optimize{
				Optimize: &evochiv1.OptimizeEvent{
					TaskId:  t.ID().String(),
					Epoch:   int32(t.Epoch),
					Rewards: f64ToF32(t.Rewards),
				},
			},
		}

	case *task.Stop:
		return &evochiv1.SubscribeResponse{
			Type: evochiv1.EventType_EVENT_TYPE_STOP,
			Event: &evochiv1.SubscribeResponse_Stop{
				Stop: &evochiv1.StopEvent{
					TaskId: t.ID().String(),
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
