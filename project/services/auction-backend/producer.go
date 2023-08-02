package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Producer struct {
	updateChannel <-chan AuctionDto
	inner         *kafka.Producer
	auctionTopic  string
	productTopic  string
	bidTopic      string
	channel       chan kafka.Event
}

func NewProducer(cfg Config, updateChannel <-chan AuctionDto) Producer {
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

	producer := Producer{
		updateChannel: updateChannel,
		inner:         p,
		auctionTopic:  cfg.Kafka.Topics.Auction,
		productTopic:  cfg.Kafka.Topics.Product,
		bidTopic:      cfg.Kafka.Topics.Bid,
		channel:       deliveryChan,
	}
	go producer.monitor()

	return producer
}

func (p *Producer) monitor() {
	for {
		update := <-p.updateChannel
		err := p.SendAuctionUpdate(update)
		if err != nil {
			fmt.Printf("error sending auction update: %s", err)
		}
		if update.Product.State != StateSubmitted {
			err = p.SendProductUpdate(update.Product)
			if err != nil {
				fmt.Printf("error sending Product update: %s", err)
			}
		}
	}
}

func (p *Producer) SendProductUpdate(dto ProductDto) error {

	message, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = p.inner.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.productTopic, Partition: kafka.PartitionAny},
			Value:          message,
			Key:            []byte(dto.ProductId),
		},
		p.channel,
	)

	return err
}

func (p *Producer) SendAuctionUpdate(dto AuctionDto) error {

	message, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = p.inner.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.auctionTopic, Partition: kafka.PartitionAny},
			Value:          message,
			Key:            []byte(dto.Product.ProductId),
		},
		p.channel,
	)

	return err
}

func (p *Producer) Close() {
	p.inner.Flush(1000)
	p.inner.Close()
}
