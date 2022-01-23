package posts

import (
	"context"
	"encoding/json"
	"sync"

	"gbu-telegram-bot/internal/bot"
	"gbu-telegram-bot/internal/entity"

	"gbu-telegram-bot/pkg/logger"
	"gbu-telegram-bot/pkg/sleep"
	"gbu-telegram-bot/pkg/wrappers/rabbit"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Posts is implementation for bot.Posts interface.
type Posts struct {
	rabbitConfig RabbitConfig
	log          logger.Logger

	rabbit *amqp.Channel
	mu     *sync.Mutex
}

var _ bot.Posts = &Posts{}

// New returns bot.Posts implementation.
func New(rabbitConfig RabbitConfig, log logger.Logger) *Posts {
	return &Posts{
		rabbitConfig: rabbitConfig,
		log:          log,

		rabbit: nil, // Initialized in Init method.
		mu:     &sync.Mutex{},
	}
}

// Init connects to rabbit and gets rabbit channel, after what
// initializes rabbit's entiies like exchanges, queues etc.
// It also registers a handler for channel closed event to reconnect.
// Close handler uses processCtx for it's calls because ctx for Init's call
// can be another: for example, limited as WithTimeout.
func (p *Posts) Init(ctx, processCtx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg := p.rabbitConfig

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

	// Exchange is fanout so no binding key required.
	err = ch.QueueBind(postsQueue, "", postsExchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "can't bind posts queue to posts exchange")
	}

	errs := make(chan *amqp.Error)
	ch.NotifyClose(errs)

	handleChannelClose := func() {
		closeErr := <-errs // This chan will get a value when rabbit channel will be closed.

		p.log.Error(errors.Wrap(closeErr, "rabbit channel closed"))

		if !conn.IsClosed() {
			err := conn.Close()
			if err != nil {
				p.log.Error(errors.Wrap(err, "can't close rabbit connection"))
			}
		}

		for attempt, isConnected := 1, false; !isConnected; attempt++ {
			sleep.WithContext(processCtx, cfg.ReconnectDelay)

			err := p.Init(processCtx, processCtx)
			if err != nil {
				p.log.Warn(errors.Wrapf(err, "can't re-init consuemr (attempt #%d)", attempt))
				continue
			}

			isConnected = true
		}

		p.log.Info("reconnected to rabbit")
	}
	go handleChannelClose()

	p.rabbit = ch

	return nil
}

func (p *Posts) Consume(ctx context.Context) (<-chan entity.PostEvent, error) {
	messages, err := p.rabbit.Consume(postsQueue, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't consume messages from queue")
	}

	posts := make(chan entity.PostEvent)

	// waitReconnection tryies to consume from queue.
	// returns false if context closed before reconnected, true otherwise.
	waitReconnection := func() bool {
		// Loop until connection reestablished or context closed.
		for {
			isCtxClosed := sleep.WithContext(ctx, p.rabbitConfig.ReconnectDelay)
			if isCtxClosed {
				return false
			}

			// TODO: Guess it can be a data race with c.rabbit.
			messages, err = p.rabbit.Consume(postsQueue, consumerName, false, false, false, false, nil)
			if err == nil {
				return true
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
				// ok is false if messages chan is closed and reconnection needed.
				if !ok {
					isReconnected := waitReconnection()
					if !isReconnected {
						close(posts)
						return
					}
					continue
				}

				err := handleMessage(message)
				if err != nil {
					p.log.Error(errors.Wrap(err, "can't handle message"))
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
