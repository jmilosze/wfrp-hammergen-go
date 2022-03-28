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

const SignalBufferLen = 2

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
	server.Start()

	c := make(chan os.Signal, SignalBufferLen)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signalRcv := <-c

	log.Printf("stop signal received: %s, starting shutdown", signalRcv)
	server.Stop()

	return nil
}
