package goem

import "github.com/urfave/cli"

var commands = []cli.Command{
	list,
}

var list = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List mails",
	Action:  handleList,
	Before:  loadListConfig,
	Flags:   listFlags,
}

var listFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "maildir, m",
		Usage: "Load Maildir from `DIR`",
	},
}
