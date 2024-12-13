package main

import (
	"context"
	"log/slog"
)

func main() {
	run(context.Background())

	slog.Debug("main debug msg") // want `context not passed to slog call`

	slog.DebugContext(context.Background(), "main debug msg with ctx")
}

func run(ctx context.Context) {
	slog.Debug("debug msg") // want `context not passed to slog call`
	slog.Info("info msg")   // want `context not passed to slog call`
	slog.Warn("warn msg")   // want `context not passed to slog call`
	slog.Error("error msg") // want `context not passed to slog call`

	slog.DebugContext(ctx, "debug msg with ctx")
	slog.InfoContext(ctx, "info msg with ctx")
	slog.WarnContext(ctx, "warn msg with ctx")
	slog.ErrorContext(ctx, "error msg with ctx")
}
