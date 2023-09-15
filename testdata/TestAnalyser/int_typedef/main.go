package main

import "log/slog"

func main() {
	slog.Debug("Hello, world", IntTypeDef(13))
}

type IntTypeDef int
