package main

import (
	"context"
	"log/slog"
)

func main() {
	run(context.Background())
}

func run(_ context.Context) { // want `context needed by slog call is blank`
	slog.Info("hello world") // want `context not passed to slog call`
}
