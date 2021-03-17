package pkg

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type MqProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	config *MqConfig
}

func NewMqProducer(config *MqConfig) (*MqProducer, error) {
	mq := &MqProducer{
		config: config,
	}
	err := mq.init()
	return mq, err
}

func (mq *MqProducer) Destroy() {
	if mq.conn.IsClosed() {
		return
	}
	_ = mq.channel.Close()
	_ = mq.conn.Close()
}

// init exchange、queue、queue bind 都做了冗余的声明操作，为了防止发送的消息
// 在 mq server 里匹配不到对应的 queue
func (mq *MqProducer) init() (err error) {
	// 建立 tcp 连接
	mq.conn, err = amqp.Dial(mq.config.Addr)
	if err != nil {
		return
	}
	// 获得 channel
	if mq.channel, err = mq.conn.Channel(); err != nil {
		return
	}
	// 声明 exchange
	if err = mq.channel.ExchangeDeclare(mq.config.Exchange, mq.config.ExchangeType, true, false, false, false, nil); err != nil {
		return
	}
	// 声明 queue
	if _, err = mq.channel.QueueDeclare(mq.config.Queue, true, false, false, false, amqp.Table{"x-max-priority": 10}); err != nil {
		return
	}
	// 绑定 queue to exchange
	if err = mq.channel.QueueBind(mq.config.Queue, mq.config.RoutingKey, mq.config.Exchange, false, nil); err != nil {
		return
	}
	return
}

func (mq *MqProducer) Publish(msg interface{}) (err error) {
	if err = mq.init(); err != nil {
		return
	}
	if _, err = mq.channel.QueueDeclare(mq.config.Queue, true, false, false, false, amqp.Table{"x-max-priority": 10}); err != nil {
		return
	}
	if err = mq.channel.QueueBind(mq.config.Queue, mq.config.RoutingKey, mq.config.Exchange, false, nil); err != nil {
		return
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	return mq.channel.Publish(
		mq.config.Exchange,
		mq.config.RoutingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
