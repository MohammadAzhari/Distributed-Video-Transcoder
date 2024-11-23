package consumer

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/communicator"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/handler"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

const (
	groupID = "my-group-1"
)

func NewConsumer(kafkaHost string, topic string, communicator *communicator.Communicator) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	broker := kafkaHost

	var group sarama.ConsumerGroup
	var err error

	for i := 0; i < 10; i++ {
		group, err = sarama.NewConsumerGroup([]string{broker}, groupID, config)
		if err == nil {
			break
		}
		log.Print("Could not connect to Kafka: ", err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	go func() {
		for {
			topics := []string{topic}
			handler := &ConsumerGroupHandler{
				handler: handler.NewHandler(communicator),
			}

			group.Consume(ctx, topics, handler)
		}
	}()

	return &Consumer{
		group: group,
	}
}

func (c *Consumer) Close() {
	c.group.Close()
}
