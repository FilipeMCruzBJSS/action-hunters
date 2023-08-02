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
	store        *Store
}

func NewConsumer(cfg Config, store *Store) *Consumer {
	kafkaConfig := cfg.ToKafkaConfig()

	c, err := kafka.NewConsumer(&kafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	consumer := &Consumer{
		bidTopic:     cfg.Kafka.Topics.Bid,
		productTopic: cfg.Kafka.Topics.Product,
		inner:        c,
		store:        store,
	}

	return consumer
}

func (c *Consumer) Subscribe() {
	err := c.inner.SubscribeTopics([]string{c.bidTopic, c.productTopic}, nil)
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
			dto, err := c.process(ev)
			if err != nil {
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

func (c *Consumer) process(event *kafka.Message) (interface{}, error) {
	switch *event.TopicPartition.Topic {
	case c.bidTopic:
		var dto BidDto
		err := json.Unmarshal(event.Value, &dto)
		if err != nil {
			fmt.Printf("Failed to consume message: %s\n", err)
			return dto, err
		}
		return dto, c.store.TryProcess(dto)
	case c.productTopic:
		var dto ProductDto
		err := json.Unmarshal(event.Value, &dto)
		if err != nil {
			fmt.Printf("Failed to consume message: %s\n", err)
			return dto, err
		}
		return dto, c.store.Store(dto)
	default:
		fmt.Printf("Failed to consume message: unknown topic")
		return nil, errors.New("unknown topic")
	}
}
