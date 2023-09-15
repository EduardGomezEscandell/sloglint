package main

import (
	"net/http"

	"log/slog"
)

func main() {
	var c http.Client
	slog.Info("This call is ill-formed", c)
}
