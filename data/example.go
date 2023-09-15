package main

import (
	"log/slog"
	"net/http"
)

func main() {

	// Not yet supported: false positive
	slog.Info("Hello", slog.Bool("Success", true))

	slog.Info("This call is ill-formed", "I'm missing my arg :(")

	slog.Warn("This call is ill-formed", 13, "bad!")

	slog.Info("This call is well-formed", getString(), "Good!")

	// This is a false negative :(
	slog.Debug("This call is ill-formed", http.Client{}, "good!")
}

func getString() string {
	return "hello, world"
}
