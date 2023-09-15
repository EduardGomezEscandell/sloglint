package main

import "log/slog"

func main() {
	slog.Info("This call is ill-formed", GetInt())
}

func GetInt() int {
	return 42
}
