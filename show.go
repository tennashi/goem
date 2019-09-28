package goem

import (
	"fmt"

	"github.com/urfave/cli"
)

func handleShow(c *cli.Context) error {
	fmt.Println(c.GlobalString("maildir"))
	return nil
}
