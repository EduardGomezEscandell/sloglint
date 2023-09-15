package main

import "log/slog"

func main() {
	slog.Debug("Hello, world", StringTypeDef("Fourty-two"))
}

type StringTypeDef string
