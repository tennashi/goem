package goem

import (
	"fmt"
	"io"
	"os"

	"github.com/tennashi/goem/shellpath"
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
	app.Before = setConfig

	return app.Run(os.Args)
}

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	cli.StringFlag{
		Name:  "maildir, m",
		Usage: "Load Maildir from `DIR`",
	},
}

const UsageText = `Usage: goem`

func setConfig(c *cli.Context) error {
	cfgPath := c.GlobalString("config")
	cfg, err := loadConfig(shellpath.Resolve(cfgPath))
	if err != nil {
		fmt.Println("config error: ", err)
		return nil
	}
	if !c.GlobalIsSet("maildir") {
		c.GlobalSet("maildir", cfg.Maildir)
	}
	return nil
}
