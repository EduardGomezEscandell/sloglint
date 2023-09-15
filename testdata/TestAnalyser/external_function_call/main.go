package main

import (
	"log"
	"log/slog"
)

func main() {
	slog.Info("This call is ill-formed", log.Flags())
}
