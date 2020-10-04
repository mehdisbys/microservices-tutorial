package common

import (
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Consumer interface {
	Receive(topic string) error
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	handler  MessageHandler
}

type MessageHandler interface {
	HandleMessage(message []byte) error
}

func NewKafkaConsumer(consumer sarama.Consumer, handler MessageHandler) *KafkaConsumer {
	return &KafkaConsumer{consumer: consumer, handler: handler}
}

func (k KafkaConsumer) Receive(topic string) error {
	defer func() {
		if err := k.consumer.Close(); err != nil {
			log.Error().Err(err).Msg("error closing consumer")
		}
	}()

	consumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error().Err(err).Msg("consumer error")
			case msg := <-consumer.Messages():
				err := k.handler.HandleMessage(msg.Value)
				if err != nil {
					log.Error().Err(err).Msg("error receiving msg")
				}
			case <-signals:
				log.Info().Msg("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	return nil
}
