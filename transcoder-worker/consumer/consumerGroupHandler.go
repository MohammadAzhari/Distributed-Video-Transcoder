package consumer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/handler"
)

type ConsumerGroupHandler struct {
	handler *handler.Handler
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Mark the beginning of a new session. This is called when the consumer group is being rebalanced.
	log.Println("Consumer group is being rebalanced")
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Mark the end of the current session. This is called just before the next rebalance happens.
	log.Println("Rebalancing will happen soon, current session will end")
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// This is where you put your message handling logic
	for msg := range claim.Messages() {
		switch string(msg.Value) {
		case "new file":
			h.handler.Init(string(msg.Key))
		case "close file":
			h.handler.End(string(msg.Key))
		default:
			h.handler.Process(string(msg.Key), msg.Value)
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n, valueLength: %v", msg.Topic, msg.Partition, msg.Offset, len(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
