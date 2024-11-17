package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/handler"
)

type Consumer struct {
	group sarama.ConsumerGroup
}

func NewConsumer(kafkaHost string, topic string) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	broker := kafkaHost
	groupID := "my-group-1"

	group, err := sarama.NewConsumerGroup([]string{broker}, groupID, config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	go func() {
		for {
			topics := []string{topic}
			handler := ConsumerGroupHandler{
				handler: handler.NewHandler(),
			}

			err := group.Consume(ctx, topics, &handler)
			if err != nil {
				panic(err)
			}
		}
	}()

	return &Consumer{
		group: group,
	}
}

func (c *Consumer) Close() {
	c.group.Close()
}
