package common

import (
	"log"

	"github.com/Shopify/sarama"
)

type Sender interface {
	Send(topic, msg string) error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaSender(producer sarama.SyncProducer) *KafkaProducer {
	return &KafkaProducer{producer: producer}
}

func (k *KafkaProducer) Send(topic, msg string) error {
	kmsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := k.producer.SendMessage(kmsg)
	if err != nil {
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
