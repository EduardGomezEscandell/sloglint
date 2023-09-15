package main

import "log/slog"

func main() {
	slog.Debug("Hello, world", StringAlias("Thirteen"))
}

type StringAlias = string
