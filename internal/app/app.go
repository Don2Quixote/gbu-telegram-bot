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

	// Getting configuration
	var cfg appConfig
	err := config.Parse(&cfg)
	if err != nil {
		return errors.Wrap(err, "can't parse config")
	}

	// Getting required connections
	pool, tgBot, err := makeConnections(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "can't make connections")
	}
	defer pool.Close()

	// Making dependencies for bot
	users, consumer, messages, err := makeDependencies(ctx, cfg, pool, tgBot, log)
	if err != nil {
		return errors.Wrap(err, "can't construct dependencies")
	}

	// Constructing and launching bot
	err = bot.New(users, consumer, messages, log).Launch(ctx)
	if err != nil {
		return errors.Wrap(err, "can't launch bot")
	}

	log.Info("app finished")

	return nil
}
