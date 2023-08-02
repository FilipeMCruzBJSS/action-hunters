package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Producer struct {
	inner   *kafka.Producer
	topic   string
	channel chan kafka.Event
}

func NewProducer(cfg Config) Producer {
	kafkaConfig := cfg.ToKafkaConfig()

	p, err := kafka.NewProducer(&kafkaConfig)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	deliveryChan := make(chan kafka.Event, 10000)

	return Producer{
		inner:   p,
		topic:   cfg.Kafka.Topics.Product,
		channel: deliveryChan,
	}
}

func (p *Producer) Send(dto OutputProductDto) error {

	message, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = p.inner.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Value:          message,
			Key:            []byte(dto.ProductId),
		},
		p.channel,
	)

	return err
}

func (p *Producer) Close() {
	p.inner.Flush(1000)
	p.inner.Close()
}
