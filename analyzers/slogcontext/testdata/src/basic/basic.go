package main

import (
	"context"
	"log/slog"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	slog.Debug("debug msg") // want `context.Context not passed to slog.Debug`
	slog.Info("info msg")   // want `context.Context not passed to slog.Info`
	slog.Warn("warn msg")   // want `context.Context not passed to slog.Warn`
	slog.Error("error msg") // want `context.Context not passed to slog.Error`
}
