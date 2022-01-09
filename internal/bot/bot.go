// Package bot provies main service's logic - notifications for telegram bot's users.
package bot

import (
	"context"
	"sync"

	"gbu-telegram-bot/pkg/logger"

	"github.com/pkg/errors"
)

// Bot is struct that incapsulates business-logic's dependencies (interfaces).
type Bot struct {
	users    Users
	consumer Consumer
	messages Messages
	log      logger.Logger
}

// New returns new bot with main business-logic of this service.
func New(users Users, consumer Consumer, messages Messages, log logger.Logger) *Bot {
	return &Bot{
		users:    users,
		consumer: consumer,
		messages: messages,
		log:      log,
	}
}

func (b *Bot) Launch(ctx context.Context) error {
	posts, err := b.consumer.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "can't consume events about new posts")
	}

	messages, err := b.messages.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "can't consume incoming messages")
	}

	wg := &sync.WaitGroup{} // WaitGroup to wait for handlers finish their job
	handlersCount := 2      // handleNewPosts and handleMessages
	wg.Add(handlersCount)

	handleNewPosts := func() {
		defer wg.Done()
		for event := range posts {
			b.handleNewPost(ctx, event)
		}
	}
	go handleNewPosts()

	handleMessages := func() {
		defer wg.Done()
		for msg := range messages {
			go b.handleMessage(ctx, msg)
		}
	}
	go handleMessages()

	<-ctx.Done()

	b.log.Info("waiting handlers to stop")
	wg.Wait()

	return nil
}
