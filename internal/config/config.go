package config

import (
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain"
	"time"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
}

type MockDbUserService struct {
	BcryptCost int
	SeedUsers  map[string]*domain.User
}

type JwtConfig struct {
	ExpiryTime time.Duration
	HmacSecret string
}

type Config struct {
	ServerConfig     *ServerConfig
	MemDbUserService *MockDbUserService
	JwtConfig        *JwtConfig
}

func NewDefault() (*Config, error) {
	users := map[string]*domain.User{
		"0": {Username: "User1", Password: "123", SharedAccounts: []string{"1"}, Admin: true},
		"1": {Username: "User2", Password: "456"},
	}

	return &Config{
		ServerConfig:     &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		MemDbUserService: &MockDbUserService{BcryptCost: 12, SeedUsers: users},
		JwtConfig:        &JwtConfig{ExpiryTime: 24 * time.Hour, HmacSecret: "some_secret"},
	}, nil
}
