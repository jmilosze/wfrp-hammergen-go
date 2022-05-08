package config

import (
	"time"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
}

type UserServiceConfig struct {
	BcryptCost int
	SeedUsers  map[string]*UserSeed
}

type JwtConfig struct {
	AccessExpiryTime time.Duration
	ResetExpiryTime  time.Duration
	HmacSecret       string
}

type EmailConfig struct {
	FromAddress string
}

type Config struct {
	ServerConfig      *ServerConfig
	UserServiceConfig *UserServiceConfig
	JwtConfig         *JwtConfig
	EmailConfig       *EmailConfig
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
			Username:       "user1@test.com",
			Password:       "123456",
			Admin:          true,
			SharedAccounts: []string{"1"},
		},
		"1": {
			Username:       "user2@test.com",
			Password:       "789123",
			Admin:          false,
			SharedAccounts: []string{},
		},
	}

	return &Config{
		ServerConfig:      &ServerConfig{Host: "localhost", Port: 8081, ShutdownTimeout: 2 * time.Second},
		UserServiceConfig: &UserServiceConfig{BcryptCost: 12, SeedUsers: users},
		JwtConfig:         &JwtConfig{AccessExpiryTime: 24 * time.Hour, ResetExpiryTime: 48 * time.Hour, HmacSecret: "some_secret"},
		EmailConfig:       &EmailConfig{FromAddress: "admin@hammergen.net"},
	}, nil
}
