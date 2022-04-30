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
	User        *domain.UserWrite
	Credentials *domain.UserWriteCredentials
	Claims      *domain.UserWriteClaims
}

func NewDefault() (*Config, error) {
	users := map[string]*UserSeed{
		"0": {
			User:        &domain.UserWrite{SharedAccounts: []string{"1"}},
			Credentials: &domain.UserWriteCredentials{Username: "user1", Password: "123"},
			Claims:      &domain.UserWriteClaims{Admin: true},
		},
		"1": {
			User:        &domain.UserWrite{SharedAccounts: []string{}},
			Credentials: &domain.UserWriteCredentials{Username: "user2", Password: "456"},
			Claims:      &domain.UserWriteClaims{Admin: false},
		},
	}

	return &Config{
		ServerConfig:     &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		MemDbUserService: &MockDbUserService{BcryptCost: 12, SeedUsers: users},
		JwtConfig:        &JwtConfig{AccessExpiryTime: 24 * time.Hour, ResetExpiryTime: 48 * time.Hour, HmacSecret: "some_secret"},
	}, nil
}
