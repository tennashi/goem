package goem

import "github.com/urfave/cli"

var commands = []cli.Command{
	list,
	show,
}

var list = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List mails",
	Action:  handleList,
}

var show = cli.Command{
	Name:    "show",
	Aliases: []string{"s"},
	Usage:   "Show mail",
	Action:  handleShow,
}
