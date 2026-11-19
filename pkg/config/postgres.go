package config

import (
	"github.com/caarlos0/env"
)

type Postgres struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	Username string `env:"POSTGRES_USERNAME,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
}

func LoadPostgres() (Postgres, error) {
	cfg := Postgres{}
	err := env.Parse(&cfg)
	return cfg, err
}
