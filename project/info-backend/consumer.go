package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Consumer struct {
	inner   *kafka.Consumer
	topic   string
	channel chan kafka.Event
}

func NewConsumer(cfg Config) Consumer {
	kafkaConfig := cfg.ToKafkaConfig()

	c, err := kafka.NewConsumer(&kafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	return Consumer{
		topic: cfg.Kafka.Consumer.Bid.Topic,
		inner: c,
	}
}

func (c *Consumer) Subscribe() {
	err := c.inner.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
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

			dto := convert(ev)

			val, err := json.Marshal(dto)
			if err != nil {
				return
			}

			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(val))
		}
	}
}

func (c *Consumer) Close() {
	_ = c.inner.Close()
}

func convert(event *kafka.Message) BidDto {
	var dto BidDto

	err := json.Unmarshal(event.Value, &dto)
	if err != nil {
		fmt.Printf("Failed to consume message: %s\n", err)
		os.Exit(1)
	}

	dto.TimeStamp = event.Timestamp.String()

	return dto
}
