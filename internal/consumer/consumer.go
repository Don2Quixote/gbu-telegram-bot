package consumer

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gbu-telegram-bot/internal/entity"

	"gbu-telegram-bot/pkg/logger"
	"gbu-telegram-bot/pkg/sleep"
	"gbu-telegram-bot/pkg/wrappers/rabbit"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Consumer is implementation for bot.Consumer interface.
type Consumer struct {
	rabbitConfig RabbitConfig
	rabbit       *amqp.Channel
	log          logger.Logger

	mu *sync.Mutex
}

// New returns bot.Consumer implementation.
func New(rabbitConfig RabbitConfig, log logger.Logger) *Consumer {
	return &Consumer{
		rabbitConfig: rabbitConfig,
		rabbit:       nil, // Initialized in Init method
		log:          log,

		mu: &sync.Mutex{},
	}
}

// Init connects to rabbit and gets rabbit channel, after what
// initializes rabbit's entiies like exchanges, queues etc.
// It also registers a handler for channel closed event to reconnect.
func (c *Consumer) Init(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cfg := c.rabbitConfig
	conn, err := rabbit.Dial(cfg.Host, cfg.User, cfg.Pass, cfg.Vhost, cfg.Amqps)
	if err != nil {
		return errors.Wrap(err, "can't connect to rabbit")
	}

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "can't get rabbit channel")
	}

	err = ch.ExchangeDeclare(postsExchange, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't declare posts exchange")
	}

	_, err = ch.QueueDeclare(postsQueue, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't declare posts queue")
	}

	// Exchange is fanout so no binding key required
	err = ch.QueueBind(postsQueue, "", postsExchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't bind posts queue to posts exchange")
	}

	errs := make(chan *amqp.Error)
	ch.NotifyClose(errs)

	handleChannelClose := func() {
		closeErr := <-errs // This chan will get a value when rabbit channel will be closed

		c.log.Error(errors.Wrap(closeErr, "rabbit channel closed"))

		if !conn.IsClosed() {
			err := conn.Close()
			if err != nil {
				c.log.Error(errors.Wrap(err, "can't close rabbit connection"))
			}
		}

		for attempt, isConnected := 1, false; !isConnected; attempt++ {
			time.Sleep(cfg.ReconnectDelay)

			err := c.Init(ctx)
			if err != nil {
				c.log.Warn(errors.Wrapf(err, "can't re-init consuemr (attempt #%d)", attempt))
				continue
			}

			isConnected = true
		}

		c.log.Info("reconnected to rabbit")
	}
	go handleChannelClose()

	c.rabbit = ch

	return nil
}

func (c *Consumer) Consume(ctx context.Context) (<-chan entity.PostEvent, error) {
	messages, err := c.rabbit.Consume(postsQueue, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't consume messages from queue")
	}

	posts := make(chan entity.PostEvent)

	waitReconnection := func() {
		// Loop until connection reestablished or context closed
		for {
			isCtxClosed := sleep.WithContext(ctx, c.rabbitConfig.ReconnectDelay)
			if isCtxClosed {
				close(posts)
				return
			}

			messages, err = c.rabbit.Consume(postsQueue, consumerName, false, false, false, false, nil)
			if err == nil {
				return
			}
		}
	}

	handleMessage := func(message amqp.Delivery) error {
		var post entity.Post
		err := json.Unmarshal(message.Body, &post)
		if err != nil {
			return errors.Wrapf(err, "can't decode message %q", message.Body)
		}

		posts <- entity.PostEvent{
			Post: post,
			Ack:  func() error { return message.Ack(false) },
			Nack: func(requeue bool) error { return message.Nack(false, requeue) },
		}

		return nil
	}

	handleMessages := func() {
		for {
			select {
			case message, ok := <-messages:
				// ok is false if messages chan is clsed and reconnection needed
				if !ok {
					waitReconnection()
					continue
				}

				err := handleMessage(message)
				if err != nil {
					c.log.Error(errors.Wrap(err, "can't handle message"))
				}
			case <-ctx.Done():
				close(posts)
				return
			}
		}
	}
	go handleMessages()

	return posts, nil
}
