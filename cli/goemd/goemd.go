package goemd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/tennashi/goem"
	"github.com/tennashi/goem/server"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, args []string, outStream, errStream io.Writer) int {
	cfgFlag := flag.String("c", "", "config path")
	flag.Parse()
	config := goem.NewConfig(*cfgFlag)

	log.SetPrefix("[goemd] ")
	log.SetOutput(errStream)

	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)
	eg.Go(func() error {
		return server.Run(ctx, config)
	})
	eg.Go(func() error {
		return Signal(ctx)
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})

	if err := eg.Wait(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func Signal(ctx context.Context) error {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <-ctx.Done():
		log.Println("signal closing")
		signal.Reset()
		return nil
	case sig := <-c:
		return fmt.Errorf("signal received: %v", sig.String())
	}
}
