package goem

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tennashi/goem/mail"
	"github.com/tennashi/goem/maildir"
	"github.com/tennashi/goem/shellpath"
	"github.com/urfave/cli"
)

func loadListConfig(c *cli.Context) error {
	cfgPath := c.GlobalString("config")
	cfg, err := loadConfig(shellpath.Resolve(cfgPath))
	if err != nil {
		fmt.Println("config error: ", err)
		return nil
	}
	c.Set("maildir", cfg.Maildir)
	return nil
}

func handleList(c *cli.Context) error {
	if !c.IsSet("maildir") {
		err := errors.New("maildir doesn't set")
		fmt.Println(err)
		return err
	}

	mdPath := shellpath.Resolve(c.String("maildir"))
	md, err := maildir.Open(mdPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	mails, err := md.Messages()
	if err != nil {
		fmt.Println(err)
		return err
	}

	var offset int
	offsetStr := c.Args().Get(0)
	if offsetStr == "" {
		offset = 0
	} else {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if offset > len(mails) {
			offset = len(mails) - 1
		}
	}

	var limit int
	limitStr := c.Args().Get(1)
	if limitStr == "" {
		limit = len(mails)
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if offset+limit > len(mails) {
		limit = len(mails) - offset
	}

	for _, m := range mails[offset : offset+limit] {
		h := mail.Header(m.Header)
		subject := h.Get("Subject")
		fmt.Println("Subject:", subject)
	}

	return nil
}
