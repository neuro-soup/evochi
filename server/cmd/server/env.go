package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	// logLevel is the log level to use.
	logLevel slog.Level = logLevelFromEnv()

	// port is the HTTP port to listen on.
	port uint = portFromEnv()
)

// logLevelFromEnv returns the log level to use from the environment.
func logLevelFromEnv() slog.Level {
	switch strings.ToLower(os.Getenv("EVOCHI_LOG_LEVEL")) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// portFromEnv returns the port to use from the environment.
func portFromEnv() uint {
	portStr := os.Getenv("EVOCHI_PORT")
	if portStr == "" {
		return 8080
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}
	return uint(port)
}
