package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const appName = "Hammergen"

type Config struct {
	Server      *Server
	UserService *UserService
	Jwt         *Jwt
	Email       *Email
	MongoDb     *MongoDb
}

type Server struct {
	Host            string        `default:"localhost" split_words:"true"`
	Port            int           `default:"8080" split_words:"true"`
	ShutdownTimeout time.Duration `default:"10s" split_words:"true"`
	RequestTimeout  time.Duration `default:"10s" split_words:"true"`
}

type UserService struct {
	BcryptCost      int      `default:"12" split_words:"true"`
	FrontEndUrl     *url.URL `default:"http://localhost:8080" split_words:"true"`
	CreateMockUsers bool     `default:"true" split_words:"true"`
}

type Jwt struct {
	AccessExpiry time.Duration `default:"24h" split_words:"true"`
	ResetExpiry  time.Duration `default:"48h" split_words:"true"`
	HmacSecret   string        `default:"some secret" split_words:"true"`
}

type Email struct {
	FromAddress   string `default:"admin@hammergen.net" split_words:"true"`
	FromName      string `default:"Hammergen Admin" split_words:"true"`
	PublicApiKey  string `default:"email public api key placeholder" split_words:"true"`
	PrivateApiKey string `default:"email private api key placeholder" split_words:"true"`
}

type MongoDb struct {
	Uri            string `default:"mongodb uri placeholder" split_words:"true"`
	DbName         string `default:"hammergenGo" split_words:"true"`
	UserCollection string `default:"user" split_words:"true"`
	CreateIndexes  bool   `default:"true" split_words:"true"`
}

func NewConfig() Config {
	var result Config

	if err := envconfig.Process(strings.ToUpper(appName), &result); err != nil {
		panic(fmt.Errorf("settings: cannot initalize app: %w", err))
	}

	return result
}
