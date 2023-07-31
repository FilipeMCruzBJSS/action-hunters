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
			Servers string `yml:"bootstrapservers" env:"KAFKA_BOOTSTRAP_SERVERS" env-default:"broker:9092"`
		} `yml:"bootstrap"`
		Producer struct {
			Topic    string `yml:"producertopic" env:"PRODUCER_TOPIC" env-default:"bids"`
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
