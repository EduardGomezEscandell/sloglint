package main

import (
	"log/slog"
)

func main() {
	var c MyType
	slog.Info("This call is ill-formed", c.value)
}

type MyType struct {
	value int
}
