package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
)

type Server struct {
	server *http.Server
	router *http.ServeMux
}

func NewServer(cfg *config.ServerConfig) *Server {
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	return &Server{
		server: server,
		router: http.NewServeMux(),
	}
}

func (s *Server) Start() {
	startedChannel := make(chan struct{})

	go func() {
		log.Printf("server starting on %s", s.server.Addr)
		startedChannel <- struct{}{}
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-startedChannel
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return s.server.Shutdown(ctx)
}
