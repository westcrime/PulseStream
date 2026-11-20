package config

import "github.com/caarlos0/env"

type Http struct {
	Port string `env:"HTTP_PORT,required"`
	Host string `env:"HTTP_HOST,required"`
}

func LoadHttp() (Http, error) {
	cfg := Http{}
	err := env.Parse(&cfg)
	return cfg, err
}
