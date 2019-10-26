package main

import (
	"os"

	"github.com/tennashi/goem/cli/goemd"
)

func main() {
	os.Exit(goemd.Run(os.Args))
}
