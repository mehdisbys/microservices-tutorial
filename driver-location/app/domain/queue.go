package domain

import (
	"errors"

	"github.com/nsqio/go-nsq"
)

// Queue is an interface to a process messages from a queue
type Queue interface {
	Process(addr string) error
}

// Message is type that represents the message included inside a queue event
type Message struct {
	Body       []byte            `json:"body"`
	Parameters map[string]string `json:"parameters"`
}

// NSQQueue holds the dependencies for the queues and implements Queue
type NSQQueue struct {
	client  *nsq.Consumer
	handler nsq.Handler
}

// NewNSQQueue creates a new NSQQueue
func NewNSQQueue(topic string, channel string, handler nsq.Handler) (*NSQQueue, error) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())

	if err != nil {
		return nil, err
	}
	return &NSQQueue{
		client:  consumer,
		handler: handler,
	}, nil
}

// Process will poll messages from the NSQ queue and and call q.handler
func (q *NSQQueue) Process(addr string) error {
	if q.handler == nil {
		return errors.New("handler not set")
	}

	q.client.AddConcurrentHandlers(q.handler, 1)

	err := q.client.ConnectToNSQD(addr)
	if err != nil {
		return err
	}
	return nil
}
