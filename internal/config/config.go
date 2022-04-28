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
	ExpiryTime time.Duration
	HmacSecret string
}

type Config struct {
	ServerConfig     *ServerConfig
	MemDbUserService *MockDbUserService
	JwtConfig        *JwtConfig
}

type UserSeed struct {
	User        *domain.User
	Credentials *domain.UserCredentials
	Claims      *domain.UserClaims
}

func NewDefault() (*Config, error) {
	users := map[string]*UserSeed{
		"0": {
			User:        &domain.User{SharedAccounts: []string{"1"}},
			Credentials: &domain.UserCredentials{Username: "User1", Password: "123"},
			Claims:      &domain.UserClaims{Admin: true},
		},
		"1": {
			User:        &domain.User{SharedAccounts: []string{}},
			Credentials: &domain.UserCredentials{Username: "User2", Password: "456"},
			Claims:      &domain.UserClaims{Admin: false},
		},
	}

	return &Config{
		ServerConfig:     &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		MemDbUserService: &MockDbUserService{BcryptCost: 12, SeedUsers: users},
		JwtConfig:        &JwtConfig{ExpiryTime: 24 * time.Hour, HmacSecret: "some_secret"},
	}, nil
}
