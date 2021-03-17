package mq

import (
	"context"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"time"
)

type RabbitMqConsumer interface {
	Start()
	Stop()
}

// ErrShouldDrop 如果接收到的消息 consumer 无法处理，希望从队列中删除，
// 需要返回这个错误
var ErrShouldDrop = errors.New("unprocessed message")

// ConsumerWorker 处理从 MQ 得到的消息
type ConsumerWorker interface {
	Consume(context.Context, []byte) error
}

// MQConsumer mq consumer 对象
type Consumer struct {
	*mq

	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}

	ctx    context.Context
	worker ConsumerWorker
}

// NewConsumer 创建一个 MQConsumer 实例
func NewConsumer(ctx context.Context, worker ConsumerWorker, config *Config) RabbitMqConsumer {
	c := &Consumer{
		mq:     newMq(config),
		ctx:    ctx,
		worker: worker,
		quit:   make(chan struct{}),
	}
	return c
}

// Start 启动 mq consumer
func (c *Consumer) Start() {
	if err := c.run(); err != nil {
		c.mq.log.Fatal().Err(err).Msg("failed to run consumer")
	}
	go c.reConnect()
	c.mq.log.Info().Msg(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan struct{})
	<-forever
}

// Stop 关闭 consumer
func (c *Consumer) Stop() {
	close(c.quit)
	c.mq.stop()
}

func (c *Consumer) run() (err error) {
	if err = c.mq.init(); err != nil {
		return
	}
	var delivery <-chan amqp.Delivery
	if delivery, err = c.channel.Consume(c.config.Queue, c.config.ConsumerTag, false, false, false, false, nil); err != nil {
		c.Stop()
		return
	}

	go c.handle(delivery)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return
}

func (c *Consumer) reConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				c.mq.log.Warn().Err(err).Msg("rabbitmq consumer - connection NotifyClose")
			}
		case err := <-c.channelNotify:
			if err != nil {
				c.mq.log.Warn().Err(err).Msg("rabbitmq consumer - channel NotifyClose")
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel(c.config.ConsumerTag, true); err != nil {
				c.mq.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
			}

			if err := c.conn.Close(); err != nil {
				c.mq.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for range c.channelNotify {
		}
		for range c.connNotify {
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				c.mq.log.Info().Msg("rabbitmq consumer - reconnect")

				if err := c.run(); err != nil {
					c.mq.log.Warn().Err(err).Msg("rabbitmq consumer - failCheck")

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}

func (c *Consumer) handle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		if err := c.worker.Consume(c.ctx, d.Body); err == nil {
			_ = d.Ack(false)
		} else {
			if errors.Is(err, ErrShouldDrop) {
				_ = d.Reject(false)
			} else {
				// 重新入队
				_ = d.Reject(true)
			}
		}
	}
}
