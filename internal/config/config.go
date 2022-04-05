package config

import (
	"time"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
}

type MockdbUserService struct {
	BcryptCost int
}

type Config struct {
	ServerConfig      *ServerConfig
	MockdbUserService *MockdbUserService
}

func NewDefault() (*Config, error) {
	return &Config{
		ServerConfig:      &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		MockdbUserService: &MockdbUserService{BcryptCost: 12},
	}, nil
}
