package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Consumer struct {
	inner        *kafka.Consumer
	bidTopic     string
	productTopic string
	auctionTopic string
	channel      chan kafka.Event
}

func NewConsumer(cfg Config) Consumer {
	kafkaConfig := cfg.ToKafkaConfig()

	c, err := kafka.NewConsumer(&kafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	return Consumer{
		bidTopic:     cfg.Kafka.Topics.Bid,
		productTopic: cfg.Kafka.Topics.Product,
		auctionTopic: cfg.Kafka.Topics.Auction,
		inner:        c,
	}
}

func (c *Consumer) Subscribe() {
	err := c.inner.SubscribeTopics([]string{c.bidTopic, c.productTopic, c.auctionTopic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.inner.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			dto, err := c.read(ev)
			if err != nil {
				fmt.Printf("Failed to consume message: %s\n", err)
				continue
			}

			val, err := json.Marshal(dto)
			if err != nil {
				continue
			}

			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(val))
		}
	}
}

func (c *Consumer) Close() {
	_ = c.inner.Close()
}

func (c *Consumer) read(event *kafka.Message) (interface{}, error) {
	switch *event.TopicPartition.Topic {
	case c.bidTopic:
		var dto BidDto
		err := json.Unmarshal(event.Value, &dto)
		return dto, err
	case c.productTopic:
		var dto ProductDto
		err := json.Unmarshal(event.Value, &dto)
		return dto, err
	case c.auctionTopic:
		var dto AuctionDto
		err := json.Unmarshal(event.Value, &dto)
		return dto, err
	default:
		fmt.Printf("Failed to consume message: unknown topic")
		return nil, errors.New("unknown topic")
	}
}
