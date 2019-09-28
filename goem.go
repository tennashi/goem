package goem

import (
	"io"
	"os"

	"github.com/urfave/cli"
)

type Goem struct {
	out    io.Writer
	errOut io.Writer
}

func NewGoem(out, errOut io.Writer) *Goem {
	return &Goem{
		out:    out,
		errOut: errOut,
	}
}

func (g *Goem) Run() int {
	if err := g.run(); err != nil {
		return 1
	}
	return 0
}

func (g *Goem) run() error {
	app := cli.NewApp()
	app.Name = "goem"
	app.Usage = UsageText
	app.Author = "tennashi"
	app.Email = "yuya.gt@gmail.com"
	app.Writer = g.out
	app.ErrWriter = g.errOut
	app.Commands = commands
	app.Flags = globalFlags

	return app.Run(os.Args)
}

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
}

const UsageText = `Usage: goem`
