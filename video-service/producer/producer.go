package producer

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	topic string
	conn  sarama.SyncProducer
}

type Message struct {
	Key   string
	Value string
}

func NewProducer(kafkaHost string, topic string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Metadata.AllowAutoTopicCreation = false

	var conn sarama.SyncProducer
	var err error
	for i := 0; i < 10; i++ {
		conn, err = sarama.NewSyncProducer([]string{kafkaHost}, config)
		if err == nil {
			break
		}
		log.Print("Could not connect to Kafka: ", err)
		time.Sleep(time.Duration(i) * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to Kafka: ", err)
	}

	log.Printf("Connected to Kafka at: %s, topic: %s", kafkaHost, topic)

	return &Producer{
		conn:  conn,
		topic: topic,
	}
}

func (p *Producer) SendMessage(message *Message) (partition int32, offset int64, err error) {
	partition, offset, err = p.conn.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(message.Value),
		Key:   sarama.StringEncoder(message.Key),
	})
	return
}

func (p *Producer) Close() error {
	return p.conn.Close()
}
