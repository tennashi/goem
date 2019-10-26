package goemd

import (
	"flag"
	"log"

	"github.com/tennashi/goem"
	"github.com/tennashi/goem/server"
)

func Run(args []string) int {
	cfgFlag := flag.String("c", "", "config path")
	flag.Parse()
	config := goem.NewConfig(*cfgFlag)

	log.SetPrefix("[goemd] ")

	return server.Run(config)
}
