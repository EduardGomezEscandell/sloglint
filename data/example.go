package main

import (
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("This call is ill-formed", "bad!")

	var c http.Client
	slog.Warn("Hello, world", c, "bad!")

	var t MyType
	slog.Debug("Hello, world", t.field, "bad!!!")

	slog.Info("Hello, world", getString(), "bad!")

	slog.Debug("Hello, world", 3*5, "bad!")

	slog.Info("This call is well-formed", "status", "good!")

	slog.Debug("This call is well-formed", http.Client{}, "good!")
}

func getString() string {
	return "hello, world"
}

type MyType struct {
	field int
}
