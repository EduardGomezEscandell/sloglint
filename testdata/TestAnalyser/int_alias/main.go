package main

import "log/slog"

func main() {
	slog.Debug("Hello, world", IntAlias(13))
}

type IntAlias = int
