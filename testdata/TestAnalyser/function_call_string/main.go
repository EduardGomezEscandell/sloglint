package main

import "log/slog"

func main() {
	slog.Info("This call is ill-formed", GetString())
}

func GetString() string {
	return "Hello!"
}
