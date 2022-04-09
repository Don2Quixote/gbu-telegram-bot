package app

import (
	"context"

	"gbu-telegram-bot/internal/bot"

	"gbu-telegram-bot/pkg/config"
	"gbu-telegram-bot/pkg/logger"

	"github.com/pkg/errors"
)

// Run runs app. If returned error is not nil, program exited
// unexpectedly and non-zero code should be returned (os.Exit(1) or log.Fatal(...)).
func Run(ctx context.Context, log logger.Logger) error {
	log.Info("starting app")

	// Getting configuration.
	var cfg appConfig
	err := config.Parse(&cfg)
	if err != nil {
		return errors.Wrap(err, "parse config")
	}

	// Getting required connections.
	pool, tgBot, err := makeConnections(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "make connections")
	}
	defer pool.Close()

	// Making dependencies for bot.
	users, posts, messages, err := makeDependencies(ctx, cfg, pool, tgBot, log)
	if err != nil {
		return errors.Wrap(err, "construct dependencies")
	}

	// Constructing and launching bot.
	bot := bot.New(users, posts, messages, log)
	err = bot.Launch(ctx)
	if err != nil {
		return errors.Wrap(err, "launch bot")
	}

	log.Info("app finished")

	return nil
}
