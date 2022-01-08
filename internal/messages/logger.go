package messages

import (
	"fmt"

	"gbu-telegram-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// tgbotapi lib has very uncool logger without possibility to disable it
// At least lib provides interface and possitibilty to implement it
// I put Warns in implementation because I don't think this lib logs something happy
type tgbotapiLogger struct {
	log logger.Logger
}

var _ tgbotapi.BotLogger = &tgbotapiLogger{}

func (l *tgbotapiLogger) Println(v ...interface{}) {
	if l.log != nil {
		message := fmt.Sprintf("tgbotapi: %v", v...)
		l.log.Warn(message)
	}
}

func (l *tgbotapiLogger) Printf(format string, v ...interface{}) {
	if l.log != nil {
		formatted := fmt.Sprintf(format, v...)
		message := fmt.Sprintf("tgbotapi: %v", formatted)
		l.log.Warnf(message)
	}
}
