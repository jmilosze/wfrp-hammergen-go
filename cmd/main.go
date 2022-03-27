package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program.
type Main struct {
	// Configuration path and parsed config data.
	Config     Config
	ConfigPath string

	// HTTP server for handling HTTP communication.
	HTTPServer *http.Server
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Run(ctx context.Context) (err error) {
	// Start the HTTP server.
	if err := m.HTTPServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// Config represents the CLI configuration file.
type Config struct {
	DB struct {
		URI string `toml:"uri"`
	} `toml:"db"`
}

// DefaultConfig returns a new instance of Config with defaults set.
func DefaultConfig() Config {
	var config Config
	config.DB.URI = DefaultURI
	return config
}

const (
	// DefaultConfigPath is the default path to the application configuration.
	DefaultConfigPath = "hammergen.conf"

	// DefaultURI is the default datasource name.
	DefaultURI = "zxc"
)
