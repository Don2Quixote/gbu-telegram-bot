package messages

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gbu-telegram-bot/internal/entity"

	"gbu-telegram-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// Messages is implementation for bot.Messages interface
type Messages struct {
	bot *tgbotapi.BotAPI
	cfg Config
	log logger.Logger

	// sending locked each sending for duration specified as cfg.SendingDelay
	sending *sync.Mutex
}

// New returns bot.Messages implementation
func New(bot *tgbotapi.BotAPI, cfg Config, log logger.Logger) *Messages {
	// "github.com/go-telegram-bot-api/telegram-bot-api/v5" has disgusting moment
	// Let me introduce it:
	_ = tgbotapi.SetLogger(&tgbotapiLogger{log: log})

	return &Messages{
		bot: bot,
		cfg: cfg,
		log: log,

		sending: &sync.Mutex{},
	}
}

func (m *Messages) Consume(ctx context.Context) (<-chan entity.IncomingMessage, error) {
	if m.cfg.Webhook != nil {
		// The laziest person ever found
		return nil, errors.New("getting updates with webhooks not implemented") // TODO: Implement
	}

	cfg := tgbotapi.NewUpdate(0)
	cfg.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage}

	messages := make(chan entity.IncomingMessage)

	handleUpdates := func() {
		updates := m.bot.GetUpdatesChan(cfg)
		for {
			select {
			case update, ok := <-updates:
				if !ok {
					close(messages)
					m.log.Error("updates chan closed")
					return
				}
				messages <- entity.IncomingMessage{
					From: entity.MessageSender{
						ID:       update.Message.From.ID,
						Username: update.Message.From.UserName,
						Name:     fmt.Sprintf("%s %s", update.Message.From.FirstName, update.Message.From.LastName),
					},
					Text: update.Message.Text,
				}
			case <-ctx.Done():
				close(messages)
				return
			}
		}
	}
	go handleUpdates()

	return messages, nil
}

// Send sends message with MarkdownV2 mode
func (m *Messages) Send(ctx context.Context, id int64, message entity.OutgoingMessage) error {
	m.sending.Lock()
	defer func() {
		time.Sleep(m.cfg.SendingDelay)
		m.sending.Unlock()
	}()

	tgMessage := tgbotapi.NewMessage(id, message.Text)
	tgMessage.ParseMode = tgbotapi.ModeMarkdownV2

	switch {
	case message.Keyboard != nil:
		tgMessage.ReplyMarkup = buildTelegramKeyboard(*message.Keyboard)
	case message.InlineKeyboard != nil:
		tgMessage.ReplyMarkup = buildTelegramInlineKeyboard(*message.InlineKeyboard)
	}

	_, err := m.bot.Send(tgMessage)
	if err != nil {
		return errors.Wrap(err, "can't send telegram message")
	}

	return nil
}
