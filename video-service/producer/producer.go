package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer() *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})

	if err != nil {
		panic(err)
	}

	return &Producer{producer: producer}
}

func (p *Producer) Produce(topic string, partition int32, key string, message string) {
	p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Key:            []byte(key),
		Value:          []byte(message),
	}, nil)
}
