package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Kafka struct {
		Bootstrap struct {
			Servers string `yml:"servers" env:"KAFKA_BOOTSTRAP_SERVERS"`
		} `yml:"bootstrap"`
		Producer struct {
			Topic string `yml:"topic" env:"PRODUCER_TOPIC"`
		} `yml:"producer"`
	} `yml:"kafka"`
}

func ReadConfig(configFile string) Config {
	fmt.Printf("Reading config file: %s\n", ConfigFile)

	var cfg Config

	err := cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		fmt.Printf("Failed to read properties: %s\n", err)
		os.Exit(1)
	}

	return cfg
}

func (cfg *Config) ToKafkaConfig() kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)
	m["bootstrap.servers"] = cfg.Kafka.Bootstrap.Servers

	return m
}
