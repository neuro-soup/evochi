package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/neuro-soup/evochi/server/internal/event"
	"github.com/neuro-soup/evochi/server/internal/worker"
	_ "github.com/neuro-soup/evochi/server/internal/worker"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
	"github.com/neuro-soup/evochi/server/pkg/service"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(log)

	workers := worker.NewPool()
	go workers.GarbageCollect(workerTimeout)

	events := event.NewQueue()

	handler := service.NewHandler(workers, events)

	mux := http.NewServeMux()
	mux.Handle(evochiv1connect.NewEvochiServiceHandler(handler))
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("starting server", "port", port)

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		h2c.NewHandler(mux, new(http2.Server)),
	)
	if err != nil {
		slog.Error("failed to start server", "error", err, "port", port)
	}

	slog.Info("server stopped", "port", port)
}
