package bot

import (
	"context"
	"fmt"

	"gbu-telegram-bot/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

// handleNewPost handles new post - sends notifications to subscribers
func (b *Bot) handleNewPost(ctx context.Context, event entity.PostEvent) {
	b.log.Infof("new blog post %q - %s", event.Post.Title, event.Post.URL)

	subscribed, err := b.users.GetSubscribedIDs(ctx)
	if err != nil {
		b.log.Error(errors.Wrap(err, "can't get subscribed users"))
		err := event.Nack(true)
		if err != nil {
			b.log.Error(errors.Wrap(err, "can't nack post event"))
		}
	}

	message := entity.OutgoingMessage{
		Text: newNotificationText(event.Post),
		InlineKeyboard: &entity.InlineKeyboard{
			[]entity.InlineButton{{Text: "Open", URL: event.Post.URL}},
		},
	}

	for _, id := range subscribed {
		err := b.messages.Send(ctx, id, message)
		if err != nil {
			b.log.Error(errors.Wrap(err, "can't send message"))
		}
	}

	err = event.Ack()
	if err != nil {
		b.log.Error(errors.Wrap(err, "can't ack event"))
	}
}

func newNotificationText(post entity.Post) string {
	// Bold title, italic author, regular summary
	format := "*%s*\n_%s_\n\n%s"
	title := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, post.Title)
	author := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, post.Author)
	summary := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, post.Summary)
	return fmt.Sprintf(format, title, author, summary)
}
