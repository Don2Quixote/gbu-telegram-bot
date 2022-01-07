package main

import (
	"context"

	"gbu-telegram-bot/internal/app"

	"gbu-telegram-bot/pkg/graceful"
	"gbu-telegram-bot/pkg/logger"

	"github.com/pkg/errors"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := logger.NewLogrus()

	graceful.OnShutdown(cancel)

	err := app.Run(ctx, log)
	if err != nil {
		err = errors.Wrap(err, "error running app")
		log.Fatal(err)
	}
}
