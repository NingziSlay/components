package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMqProducer interface {
	Destroy()
	Publish(interface{}) error
	PurgeQueue() error
}

type Producer struct {
	*mq
}

func NewMqProducer(config *Config) (RabbitMqProducer, error) {
	mq := newMq(config)
	if err := mq.init(); err != nil {
		return nil, err
	}
	return &Producer{
		mq: mq,
	}, nil
}

func (producer *Producer) Destroy() {
	producer.mq.stop()
}

func (producer *Producer) Publish(msg interface{}) (err error) {

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	return producer.channel.Publish(
		producer.config.Exchange,
		producer.config.RoutingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// PurgeQueue will purge all undelivered message of queue which
// declare in Config struct
func (producer *Producer) PurgeQueue() error {
	_, err := producer.channel.QueuePurge(producer.config.Queue, true)
	return err
}
