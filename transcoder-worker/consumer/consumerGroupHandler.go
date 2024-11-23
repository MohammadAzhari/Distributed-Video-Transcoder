package consumer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/MohammadAzhari/Distributed-Video-Transcoder/transcoder-worker/handler"
)

const (
	newFile   = "new file"
	closeFile = "close file"
)

type ConsumerGroupHandler struct {
	handler *handler.Handler
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer group is being rebalanced")
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Rebalancing will happen soon, current session will end")
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		key := string(msg.Key)
		log.Printf("Message claimed: key=%s topic=%q partition=%d offset=%d\n", key, msg.Topic, msg.Partition, msg.Offset)
		switch string(msg.Value) {
		case newFile:
			h.handler.Init(key)
		case closeFile:
			h.handler.End(key)
		default:
			h.handler.Process(key, msg.Value)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
