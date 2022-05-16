package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	EnvServerHost                   = "SERVER_HOST"
	EnvServerPort                   = "SERVER_PORT"
	EnvServerShutdownTimeoutSeconds = "SERVER_SHUTDOWN_TIMEOUT_SEC"
	EnvServerRequestTimeoutSeconds  = "SERVER_REQUEST_TIMEOUT_SEC"
	EnvUserBcryptCost               = "USER_BCRYPT_COST"
	EnvUserSeed                     = "USER_SEED"
	EnvJwtAccessExpirySeconds       = "JWT_ACCESS_EXPIRY_SEC"
	EnvJwtResetExpirySeconds        = "JWT_RESET_EXPIRY_SEC"
	EnvJwtHmacSecret                = "JWT_HMAC_SECRET"
	EnvEmailFromAddress             = "EMAIL_FROM_ADDRESS"
	EnvMongoDbUri                   = "MONGODB_URI"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

type UserServiceConfig struct {
	BcryptCost int
	SeedUsers  []*UserSeed
}

type JwtConfig struct {
	AccessExpiry time.Duration
	ResetExpiry  time.Duration
	HmacSecret   string
}

type EmailConfig struct {
	FromAddress string
}

type MongoDbConfig struct {
	Uri string
}

type Config struct {
	ServerConfig      *ServerConfig
	UserServiceConfig *UserServiceConfig
	JwtConfig         *JwtConfig
	EmailConfig       *EmailConfig
	MongoDbConfig     *MongoDbConfig
}

type UserSeed struct {
	Id                string
	Username          string
	Password          string
	Admin             bool
	SharedAccountsIds []string
}

func NewDefault() *Config {
	users := []*UserSeed{
		{
			Id:                "0",
			Username:          "user1@test.com",
			Password:          "123456",
			Admin:             true,
			SharedAccountsIds: []string{"1"},
		},
		{
			Id:                "1",
			Username:          "user2@test.com",
			Password:          "789123",
			Admin:             false,
			SharedAccountsIds: []string{},
		},
		{
			Id:                "2",
			Username:          "user3@test.com",
			Password:          "111111",
			Admin:             false,
			SharedAccountsIds: []string{"0", "1"},
		},
	}

	return &Config{
		ServerConfig: &ServerConfig{
			Host:            "localhost",
			Port:            8080,
			ShutdownTimeout: 10 * time.Second,
			RequestTimeout:  10 * time.Second,
		},
		UserServiceConfig: &UserServiceConfig{
			BcryptCost: 12,
			SeedUsers:  users,
		},
		JwtConfig: &JwtConfig{
			AccessExpiry: 24 * time.Hour,
			ResetExpiry:  48 * time.Hour,
			HmacSecret:   "some_secret",
		},
		EmailConfig:   &EmailConfig{FromAddress: "admin@hammergen.net"},
		MongoDbConfig: &MongoDbConfig{Uri: ""},
	}
}

func NewFromEnv() (*Config, error) {
	cfg := NewDefault()
	var err error

	cfg.ServerConfig.Host = readEnv(EnvServerHost, cfg.ServerConfig.Host)
	cfg.ServerConfig.Port, err = strconv.Atoi(readEnv(EnvServerPort, fmt.Sprintf("%d", cfg.ServerConfig.Port)))
	if err != nil {
		return nil, err
	}

	cfg.MongoDbConfig.Uri = readEnv(EnvMongoDbUri, cfg.MongoDbConfig.Uri)

	return cfg, nil
}

func readEnv(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
