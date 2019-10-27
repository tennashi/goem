package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tennashi/goem"
	"github.com/tennashi/goem/server/handler"
)

func Run(ctx context.Context, config *goem.Config) error {
	s := newServer(config)
	return s.run(ctx)
}

type server struct {
	config *goem.Config
}

func newServer(config *goem.Config) *server {
	return &server{config}
}

func (s *server) run(ctx context.Context) error {
	log.Println("server intializing")
	r := newRouter(s.config.RootDir)
	hs := &http.Server{
		Addr:    ":" + s.config.Server.Port,
		Handler: r,
	}
	log.Println("server intialized")
	log.Printf("server running on localhost:%v", s.config.Server.Port)

	eCh := make(chan error)
	go func() {
		defer close(eCh)
		if err := hs.ListenAndServe(); err != nil {
			eCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("server shuting down")
		return hs.Shutdown(ctx)
	case err := <-eCh:
		return err
	}
}

func newRouter(rootPath string) *chi.Mux {
	r := chi.NewRouter()

	h := handler.New(rootPath)
	r.Get("/maildir/", h.ListMaildir)
	r.Get("/maildir/{dirName}", h.ListMail)
	r.Get("/maildir/{dirName}/{key}", h.GetMail)

	return r
}
