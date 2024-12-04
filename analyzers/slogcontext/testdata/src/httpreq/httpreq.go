package main

import (
	"log/slog"
	"net/http"
)

func main() {
	run(&http.Request{})
}

func run(r *http.Request) {
	slog.Info("hello world") // want `context not passed to slog call`
}
