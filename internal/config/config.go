package config

import (
	"fmt"
	"net/url"
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
	EnvJwtAccessExpirySeconds       = "JWT_ACCESS_EXPIRY_SEC"
	EnvJwtResetExpirySeconds        = "JWT_RESET_EXPIRY_SEC"
	EnvJwtHmacSecret                = "JWT_HMAC_SECRET"
	EnvEmailFromAddress             = "EMAIL_FROM_ADDRESS"
	EnvEmailFromName                = "EMAIL_FROM_NANME"
	EnvEmailPublicApiKey            = "EMAIL_PUBLIC_API_KEY"
	EnvEmailPrivateApiKey           = "EMAIL_PRIVATE_API_KEY"
	EnvMongoDbUri                   = "MONGODB_URI"
)

type ServerConfig struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

type UserServiceConfig struct {
	BcryptCost  int
	FrontEndUrl *url.URL
}

type JwtConfig struct {
	AccessExpiry time.Duration
	ResetExpiry  time.Duration
	HmacSecret   string
}

type EmailConfig struct {
	FromAddress   string
	FromName      string
	PublicApiKey  string
	PrivateApiKey string
}

type MongoDbConfig struct {
	Uri            string
	DbName         string
	UserCollection string
	CreateIndexes  bool
}

type Config struct {
	ServerConfig      *ServerConfig
	UserServiceConfig *UserServiceConfig
	JwtConfig         *JwtConfig
	EmailConfig       *EmailConfig
	MongoDbConfig     *MongoDbConfig
}

type UserSeed struct {
	Id             string
	Username       string
	Password       string
	Admin          bool
	SharedAccounts []string
}

func NewDefault() *Config {

	frontEndUrl, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}

	return &Config{
		ServerConfig: &ServerConfig{
			Host:            "localhost",
			Port:            8080,
			ShutdownTimeout: 10 * time.Second,
			RequestTimeout:  10 * time.Second,
		},

		UserServiceConfig: &UserServiceConfig{
			BcryptCost:  12,
			FrontEndUrl: frontEndUrl,
		},
		JwtConfig: &JwtConfig{
			AccessExpiry: 24 * time.Hour,
			ResetExpiry:  48 * time.Hour,
			HmacSecret:   "some_secret",
		},
		EmailConfig:   &EmailConfig{FromAddress: "admin@hammergen.net", FromName: "Hammergen Admin"},
		MongoDbConfig: &MongoDbConfig{DbName: "hammergenGo", UserCollection: "user", CreateIndexes: true},
	}
}

func NewFromEnv() *Config {
	cfg := NewDefault()

	var err error
	cfg.ServerConfig.Host = readEnv(EnvServerHost, cfg.ServerConfig.Host)
	cfg.ServerConfig.Port, err = strconv.Atoi(readEnv(EnvServerPort, fmt.Sprintf("%d", cfg.ServerConfig.Port)))
	if err != nil {
		panic(err)
	}

	cfg.MongoDbConfig.Uri = readEnv(EnvMongoDbUri, cfg.MongoDbConfig.Uri)

	cfg.EmailConfig.FromAddress = readEnv(EnvEmailFromAddress, cfg.EmailConfig.FromAddress)
	cfg.EmailConfig.FromName = readEnv(EnvEmailFromName, cfg.EmailConfig.FromName)
	cfg.EmailConfig.PublicApiKey = readEnv(EnvEmailPublicApiKey, cfg.EmailConfig.PublicApiKey)
	cfg.EmailConfig.PrivateApiKey = readEnv(EnvEmailPrivateApiKey, cfg.EmailConfig.PrivateApiKey)

	return cfg
}

func readEnv(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func NewMockUsers() []*UserSeed {
	return []*UserSeed{
		{
			Id:             "000000000000000000000000",
			Username:       "user1@test.com",
			Password:       "123456",
			Admin:          true,
			SharedAccounts: []string{},
		},
		{
			Id:             "000000000000000000000001",
			Username:       "user2@test.com",
			Password:       "789123",
			Admin:          false,
			SharedAccounts: []string{"user1@test.com"},
		},
		{
			Id:             "000000000000000000000002",
			Username:       "user3@test.com",
			Password:       "111111",
			Admin:          false,
			SharedAccounts: []string{"user1@test.com", "user2@test.com"},
		},
	}
}
