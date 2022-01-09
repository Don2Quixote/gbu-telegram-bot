package app

import (
	"context"

	"gbu-telegram-bot/pkg/wrappers/pgxpool"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// makeConnections makes required connections/clients.
func makeConnections(ctx context.Context, cfg appConfig) (*pgxpool.Pool, *tgbotapi.BotAPI, error) {
	pool, err := pgxpool.Connect(ctx, cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDBName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't connect to postgres")
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		pool.Close()
		return nil, nil, errors.Wrap(err, "can't create bot API")
	}

	return pool, bot, nil
}
