package bot

import (
	"context"
	"strings"
	"time"

	"gbu-telegram-bot/internal/entity"

	"github.com/pkg/errors"
)

// handleCMessage handles incoming message.
// It calls one of methods with command's handling depending on message's text.
func (b *Bot) handleMessage(ctx context.Context, msg entity.IncomingMessage) {
	start := time.Now()
	command := parseCommand(msg.Text)
	switch command {
	case entity.CommandStart:
		b.handleStart(ctx, msg.From)
	case entity.CommandSubscribe:
		b.handleSubscribe(ctx, msg.From)
	case entity.CommandUnsubscribe:
		b.handleUnsubscribe(ctx, msg.From)
	default:
		b.handleUnknown(ctx, msg.From)
	}
	b.log.Infof("handled command %s from %s (%v)", command, msg.From.String(), time.Since(start))
}

func (b *Bot) handleStart(ctx context.Context, from entity.MessageSender) {
	_, err := b.users.Get(ctx, from.ID)
	if errors.Is(err, entity.ErrUserNotFound) {
		err := b.users.Add(ctx, entity.User{
			ID:           from.ID,
			Username:     from.Username,
			Name:         from.Name,
			IsSubscribed: true,
		})
		if err != nil {
			b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't add user"), from)
			return
		}
	}
	if err != nil && !errors.Is(err, entity.ErrUserNotFound) {
		b.log.Debugf("err is %v", err)
		b.log.Debugf("entity.ErrUserNotFound", entity.ErrUserNotFound)
		b.log.Debugf("err == entity.ErrUserNotFound", err == entity.ErrUserNotFound)
		b.log.Debugf("errors.Is(err, entity.ErrUserNotFound)", errors.Is(err, entity.ErrUserNotFound))
		b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't get user"), from)
		return
	}

	err = b.messages.Send(ctx, from.ID, greetingMessage)
	if err != nil {
		b.log.Error(errors.Wrapf(err, "can't send greeting message to %s", from.String()))
	}
}

func (b *Bot) handleSubscribe(ctx context.Context, from entity.MessageSender) {
	_, err := b.users.Get(ctx, from.ID)
	if err != nil {
		b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't get user"), from)
		return
	}

	err = b.users.Subscribe(ctx, from.ID)
	if err != nil {
		b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't subscribe user"), from)
		return
	}

	err = b.messages.Send(ctx, from.ID, subscribedMessage)
	if err != nil {
		b.log.Error(errors.Wrapf(err, "can't send unsubscribed message to %s", from.String()))
	}
}

func (b *Bot) handleUnsubscribe(ctx context.Context, from entity.MessageSender) {
	_, err := b.users.Get(ctx, from.ID)
	if err != nil {
		b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't get user"), from)
		return
	}

	err = b.users.Unsubscribe(ctx, from.ID)
	if err != nil {
		b.replyWithErrorMessage(ctx, errors.Wrap(err, "can't unsubscribe user"), from)
		return
	}

	err = b.messages.Send(ctx, from.ID, unsubscribedMessage)
	if err != nil {
		b.log.Error(errors.Wrapf(err, "can't send unsubscribed message to %s", from.String()))
	}
}

func (b *Bot) handleUnknown(ctx context.Context, from entity.MessageSender) {
	err := b.messages.Send(ctx, from.ID, unknownMessage)
	if err != nil {
		b.log.Error(errors.Wrapf(err, "can't send unknown message to %s", from.String()))
	}
}

// replyWithErrorMessage logs error, sends errorMessage to user
// and logs error if it occurs on sending.
// As error can happen in many places so it's better to wrap this logic in this
// method to reduce code's length.
func (b *Bot) replyWithErrorMessage(ctx context.Context, err error, from entity.MessageSender) {
	b.log.Error(err)
	err = b.messages.Send(ctx, from.ID, errorMessage)
	if err != nil {
		b.log.Error(errors.Wrapf(err, "can't send error message to %s", from.String()))
	}
}

// parseCommand parses command from message's text.
func parseCommand(message string) entity.Command {
	message = strings.ToLower(message)
	switch message {
	case "/start", "start", "/restart", "restart":
		return entity.CommandStart
	case "/subscribe", "📥 subscribe", "subscribe":
		return entity.CommandSubscribe
	case "/unsubscribe", "❌ unsubscribe", "unsubscribe":
		return entity.CommandUnsubscribe
	default:
		return entity.CommandUnknown
	}
}
