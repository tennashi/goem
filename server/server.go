package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tennashi/goem"
	"github.com/tennashi/goem/server/handler"
)

func Run(config *goem.Config) int {
	s := newServer(config)
	if err := s.run(); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

type server struct {
	config *goem.Config
}

func newServer(config *goem.Config) *server {
	return &server{config}
}

func (s *server) run() error {
	log.Println("server intializing")
	r := newRouter(s.config.RootDir)
	hs := &http.Server{
		Addr:    ":" + s.config.Server.Port,
		Handler: r,
	}
	log.Println("server intialized")
	log.Printf("server running on localhost:%v", s.config.Server.Port)
	if err := hs.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func newRouter(rootPath string) *chi.Mux {
	r := chi.NewRouter()

	h := handler.New(rootPath)
	r.Get("/maildir/", h.ListMaildir)
	r.Get("/maildir/{dirName}", h.ListMail)
	r.Get("/maildir/{dirName}/{key}", h.GetMail)

	return r
}
