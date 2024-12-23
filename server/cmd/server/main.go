package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	_ "github.com/neuro-soup/evochi/server/internal/worker"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(log)

	mux := http.NewServeMux()
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
