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
			Servers string `yml:"servers" env:"KAFKA_BOOTSTRAP_SERVERS" env-default:"broker:9092"`
		} `yml:"bootstrap"`
		GroupId         string `yml:"groupId" env:"KAFKA_GROUPID" env-default:"info-backend"`
		AutoOffsetReset string `yml:"offsetReset" env:"KAFKA_AUTO_OFFSET_RESET" env-default:"earliest"`
		Consumer        struct {
			Bid struct {
				Topic string `yml:"topic" env:"CONSUMER_BID_TOPIC" env-default:"bids"`
			} `yml:"bid"`
			Product struct {
				Topic string `yml:"topic" env:"CONSUMER_PRODUCT_TOPIC" env-default:"products"`
			} `yml:"product"`
		} `yml:"consumer"`
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
	m["group.id"] = cfg.Kafka.GroupId
	m["auto.offset.reset"] = cfg.Kafka.AutoOffsetReset

	return m
}
