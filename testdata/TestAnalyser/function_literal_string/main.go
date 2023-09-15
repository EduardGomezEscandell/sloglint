package main

import "log/slog"

func main() {
	slog.Info("This call is ill-formed", func() string {
		return "Fourty-two"
	}())
}
