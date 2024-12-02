package main

import (
	"context"
	"log/slog"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	slog.Debug("debug msg") // want `context not passed to slog call`
	slog.Info("info msg")   // want `context not passed to slog call`
	slog.Warn("warn msg")   // want `context not passed to slog call`
	slog.Error("error msg") // want `context not passed to slog call`
}
