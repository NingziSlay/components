package pkg

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"time"
)

// ErrShouldDrop 如果接收到的消息 consumer 无法处理，希望从队列中删除，
// 需要返回这个错误
var ErrShouldDrop = errors.New("unprocessed message")

// ConsumerWorker 处理从 MQ 得到的消息
type ConsumerWorker interface {
	Consume(context.Context, []byte) error
}

// MqConfig MQ 的基本配置
type MqConfig struct {
	Addr         string
	Exchange     string
	ExchangeType string // topic, direct, etc
	Queue        string
	RoutingKey   string
	ConsumerTag  string
}

// MQConsumer mq consumer 对象
type MQConsumer struct {
	ctx           context.Context
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
	config        *MqConfig
	worker        ConsumerWorker

	log zerolog.Logger
}

// NewConsumer 创建一个 MQConsumer 实例
func NewConsumer(ctx context.Context, worker ConsumerWorker, config *MqConfig, ) *MQConsumer {
	c := &MQConsumer{
		ctx:    ctx,
		worker: worker,
		quit:   make(chan struct{}),
		config: config,
		log:    GetLogger(),
	}
	return c
}

// Start 启动 mq consumer
func (c *MQConsumer) Start() {
	if err := c.run(); err != nil {
		c.log.Fatal().Err(err).Msg("failed to run consumer")
	}
	go c.reConnect()
	c.log.Info().Msg(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan struct{})
	<-forever
}

// Stop 关闭 consumer
func (c *MQConsumer) Stop() {
	close(c.quit)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.config.ConsumerTag, true); err != nil {
			c.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
		}

		if err := c.conn.Close(); err != nil {
			c.log.Warn().Err(err).Msg("rabbitmq consumer - connection close failed")
		}
	}
}

func (c *MQConsumer) run() error {
	var err error
	if c.conn, err = amqp.Dial(c.config.Addr); err != nil {
		return err
	}

	if c.channel, err = c.conn.Channel(); err != nil {
		c.Stop()
		return err
	}

	if _, err = c.channel.QueueDeclare(c.config.Queue, true, false, false, false, amqp.Table{"x-max-priority": 10}); err != nil {
		c.Stop()
		return err
	}

	_ = c.channel.Qos(1, 0, true)

	if err = c.channel.QueueBind(c.config.Queue, c.config.RoutingKey, c.config.Exchange, false, nil); err != nil {
		c.Stop()
		return err
	}

	var delivery <-chan amqp.Delivery
	if delivery, err = c.channel.Consume(c.config.Queue, c.config.ConsumerTag, false, false, false, false, nil); err != nil {
		c.Stop()
		return err
	}

	go c.handle(delivery)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return err
}

func (c *MQConsumer) reConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				c.log.Warn().Err(err).Msg("rabbitmq consumer - connection NotifyClose")
			}
		case err := <-c.channelNotify:
			if err != nil {
				c.log.Warn().Err(err).Msg("rabbitmq consumer - channel NotifyClose")
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// close message delivery
			if err := c.channel.Cancel(c.config.ConsumerTag, true); err != nil {
				c.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
			}

			if err := c.conn.Close(); err != nil {
				c.log.Warn().Err(err).Msg("rabbitmq consumer - channel cancel failed")
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				c.log.Info().Msg("rabbitmq consumer - reconnect")

				if err := c.run(); err != nil {
					c.log.Warn().Err(err).Msg("rabbitmq consumer - failCheck")

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}

func (c *MQConsumer) handle(delivery <-chan amqp.Delivery) {
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
