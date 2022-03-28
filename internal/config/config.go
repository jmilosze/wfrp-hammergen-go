package config

type ServerConfig struct {
	Host string
	Port int
}

type Config struct {
	APIServer *ServerConfig
}

func NewDefault() (*Config, error) {
	return &Config{
		APIServer: &ServerConfig{Host: "localhost", Port: 8080},
	}, nil
}
