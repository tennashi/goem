package main

import (
	"os"

	"github.com/tennashi/goem"
)

func main() {
	g := goem.NewGoem(os.Stdout, os.Stderr)
	os.Exit(g.Run())
}
