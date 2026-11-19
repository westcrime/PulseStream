package config

import "github.com/caarlos0/env"

type Kafka struct {
	Host    string `env:"KAFKA_HOST,required"`
	Port    string `env:"KAFKA_PORT,required"`
	Topic   string `env:"KAFKA_TOPIC,required"`
	GroupID string `env:"KAFKA_GROUP_ID,required"`
}

func LoadKafka() (Kafka, error) {
	cfg := Kafka{}
	err := env.Parse(&cfg)
	return cfg, err
}
