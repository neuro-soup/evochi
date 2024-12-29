package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/neuro-soup/evochi/server/internal/distribution/worker"
	v1 "github.com/neuro-soup/evochi/server/pkg/handler/v1"
	"github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1/evochiv1connect"
)

func main() {
	configureLogger(slog.LevelInfo)

	cfg, err := loadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	configureLogger(cfg.SlogLevel())

	workers := worker.NewPool()
	go workers.Watch(cfg.WorkerTimeout)

	mux := http.NewServeMux()
	registerV1(cfg, mux, workers)
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("starting server", "port", cfg.ServerPort)

	err = http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.ServerPort),
		h2c.NewHandler(mux, new(http2.Server)),
	)
	if err != nil {
		slog.Error("failed to start server", "error", err, "port", cfg.ServerPort)
	}
}

func configureLogger(level slog.Level) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(log)
}

func registerV1(cfg *config, mux *http.ServeMux, workers *worker.Pool) {
	slog.Debug("registering v1 handler")

	attrs := make(map[string][]byte)
	for key, val := range cfg.Attrs {
		attrs[key] = []byte(val)
	}

	v1 := v1.New(
		v1.Config{
			JWTSecret:      cfg.JWTSecret,
			MaxWorkers:     cfg.MaxWorkers,
			WorkerTimeout:  cfg.WorkerTimeout,
			MaxEpochs:      cfg.MaxEpochs,
			PopulationSize: cfg.PopulationSize,
			Attrs:          attrs,
		},
		workers,
	)

	mux.Handle(evochiv1connect.NewEvochiServiceHandler(v1))
}
