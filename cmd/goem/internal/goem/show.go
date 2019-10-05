package goem

import (
	"errors"
	"fmt"

	"github.com/tennashi/goem/maildir"
	"github.com/tennashi/goem/shellpath"
	"github.com/urfave/cli"
)

func handleShow(c *cli.Context) error {
	if !c.GlobalIsSet("maildir") {
		err := errors.New("maildir doesn't set")
		fmt.Println(err)
		return err
	}

	mdPath := shellpath.Resolve(c.GlobalString("maildir"))
	md, err := maildir.New(mdPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	key := c.Args().Get(0)
	if key == "" {
		err := errors.New("key is required")
		fmt.Println(err)
		return err
	}

	mail, err := md.GetMessageWithRawKey(key, maildir.SubDirCur)
	if err != nil {
		return err
	}

	fmt.Println(mail)

	return nil
}
