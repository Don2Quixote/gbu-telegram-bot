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
	posts    Posts
	messages Messages
	log      logger.Logger
}

// New returns new bot with main business-logic of this service.
func New(users Users, posts Posts, messages Messages, log logger.Logger) *Bot {
	return &Bot{
		users:    users,
		posts:    posts,
		messages: messages,
		log:      log,
	}
}

// Launch launches event's consuming and their's handling.
// Blocks until context closed or launching error happened.
// Returns nil error if context closed.
func (b *Bot) Launch(ctx context.Context) error {
	// posts chan will be closed when context will be closed.
	posts, err := b.posts.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "can't consume events about new posts")
	}

	// messages chan will be closed when context will be closed.
	messages, err := b.messages.Consume(ctx)
	if err != nil {
		return errors.Wrap(err, "can't consume incoming messages")
	}

	wg := &sync.WaitGroup{} // WaitGroup to wait for handlers finish their job.
	handlersCount := 2      // handleNewPosts and handleMessages.
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
