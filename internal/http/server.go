package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(cfg *config.ServerConfig, router http.Handler) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: router,
	}

	return &Server{
		server:          server,
		shutdownTimeout: cfg.ShutdownTimeout,
	}
}

func (s *Server) Start() {
	go func() {
		log.Printf("server starting on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
