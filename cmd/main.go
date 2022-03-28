package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg, err := config.NewDefault()
	if err != nil {
		return fmt.Errorf("getting service config from environment: %w", err)
	}

	server := http.NewServer(cfg.APIServer)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	server.Start()

	<-done

	if err := server.Stop(); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}
