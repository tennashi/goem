package main

import (
	"context"
	"os"

	"github.com/tennashi/goem/cli/goemd"
)

func main() {
	os.Exit(goemd.Run(
		context.Background(),
		os.Args,
		os.Stdout,
		os.Stderr,
	))
}
