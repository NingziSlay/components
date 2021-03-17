package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMqProducer interface {
	Destroy()
	Publish(interface{}) error
}

type Producer struct {
	*mq
}

func NewMqProducer(config *Config) RabbitMqProducer {
	return &Producer{
		newMq(config),
	}
}

func (producer *Producer) Destroy() {
	producer.mq.stop()
}

func (producer *Producer) Publish(msg interface{}) (err error) {
	if err = producer.mq.init(); err != nil {
		return
	}

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
