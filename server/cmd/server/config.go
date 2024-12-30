package main

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	// LogLevel is the log level to use. Either "debug", "info", "warn", or "error".
	// Defaults to "info".
	LogLevel string `split_words:"true" default:"info"`

	// ServerPort is the port to listen on. Defaults to 8080.
	ServerPort uint `split_words:"true" default:"8080"`

	// JWTSecret is the secret to use for JWT tokens.
	JWTSecret string `split_words:"true" required:"true"`

	// WorkerTimeout is the timeout for a worker. Defaults to 1 minute.
	WorkerTimeout time.Duration `split_words:"true" default:"1m"`

	// MaxWorkers is the maximum number of workers to run. Defaults to 0, which
	// means there is no limit.
	MaxWorkers uint `split_words:"true" default:"0"`

	// MaxEpochs is the maximum number of epochs to run. Defaults to 0, which
	// means there is no limit.
	MaxEpochs uint `split_words:"true" default:"0"`

	// PopulationSize is the size of the population to use.
	PopulationSize uint `split_words:"true" required:"true"`

	// Attrs are the custom attributes to use.
	Attrs map[string]string `split_words:"true"`
}

func loadConfig() (*config, error) {
	var cfg config
	if err := envconfig.Process("EVOCHI", &cfg); err != nil {
		return nil, fmt.Errorf("config: failed to load: %w", err)
	}
	return &cfg, nil
}

func (c config) SlogLevel() slog.Level {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
