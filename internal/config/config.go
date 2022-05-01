package config

import (
	"time"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
}

type MockDbUserService struct {
	BcryptCost int
	SeedUsers  map[string]*UserSeed
}

type JwtConfig struct {
	AccessExpiryTime time.Duration
	ResetExpiryTime  time.Duration
	HmacSecret       string
}

type Config struct {
	ServerConfig     *ServerConfig
	MemDbUserService *MockDbUserService
	JwtConfig        *JwtConfig
}

type UserSeed struct {
	Username       string
	Password       string
	Admin          bool
	SharedAccounts []string
}

func NewDefault() (*Config, error) {
	users := map[string]*UserSeed{
		"0": {
			Username:       "user1",
			Password:       "123",
			Admin:          true,
			SharedAccounts: []string{"1"},
		},
		"1": {
			Username:       "user2",
			Password:       "456",
			Admin:          false,
			SharedAccounts: []string{},
		},
	}

	return &Config{
		ServerConfig:     &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		MemDbUserService: &MockDbUserService{BcryptCost: 12, SeedUsers: users},
		JwtConfig:        &JwtConfig{AccessExpiryTime: 24 * time.Hour, ResetExpiryTime: 48 * time.Hour, HmacSecret: "some_secret"},
	}, nil
}
