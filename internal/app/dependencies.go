package app

import (
	"context"
	"time"

	"gbu-telegram-bot/internal/bot"
	"gbu-telegram-bot/internal/messages"
	"gbu-telegram-bot/internal/posts"
	"gbu-telegram-bot/internal/users"

	"gbu-telegram-bot/pkg/logger"
	"gbu-telegram-bot/pkg/wrappers/pgxpool"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// makeDependencies makes all bot's dependencies.
func makeDependencies(
	ctx context.Context,
	cfg appConfig,
	pool *pgxpool.Pool,
	bot *tgbotapi.BotAPI,
	log logger.Logger,
) (
	bot.Users,
	bot.Posts,
	bot.Messages,
	error,
) {
	users := users.New(pool)

	posts := posts.New(posts.RabbitConfig{
		Host:           cfg.RabbitHost,
		User:           cfg.RabbitUser,
		Pass:           cfg.RabbitPass,
		Vhost:          cfg.RabbitVhost,
		Amqps:          cfg.RabbitAmqps,
		ReconnectDelay: time.Duration(cfg.RabbitReconnectDelay) * time.Second,
	}, log)

	err := posts.Init(ctx, ctx)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can't init posts")
	}

	messages := messages.New(bot, messages.Config{
		Webhook:      nil, // TODO: Should be configurable, but webhooks not implemented
		SendingDelay: time.Millisecond * time.Duration(cfg.MessagesSendingDelay),
	}, log)

	return users, posts, messages, nil
}
