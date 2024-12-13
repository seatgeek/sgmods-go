package main

import "log/slog"

func main() {
	slog.Info("hello world!") //lint:ignore slogcontext "no context in main"
	slog.Info("hello world!") //lint:ignore otherlinter "another reason" // want `context not passed to slog call`
}
