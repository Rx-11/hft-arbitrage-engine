package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	ApiKey   string `env:"API_KEY" envDefault:"-"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to system env vars")
	}
	Cfg = &Config{}
	if err := env.Parse(Cfg); err != nil {
		return nil, err
	}
	return Cfg, nil
}
