package config

import "time"

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
}

type Config struct {
	APIServer *ServerConfig
}

func NewDefault() (*Config, error) {
	return &Config{
		APIServer: &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
	}, nil
}
