package messages

import (
	"gbu-telegram-bot/internal/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func buildTelegramKeyboard(keyboard entity.Keyboard) tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0, len(keyboard))

	for _, row := range keyboard {
		buttons := make([]tgbotapi.KeyboardButton, 0, len(row))

		for _, buttonText := range row {
			buttons = append(buttons, tgbotapi.NewKeyboardButton(buttonText))
		}

		rows = append(rows, tgbotapi.NewKeyboardButtonRow(buttons...))
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}

func buildTelegramInlineKeyboard(keyboard entity.InlineKeyboard) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(keyboard))

	for _, row := range keyboard {
		buttons := make([]tgbotapi.InlineKeyboardButton, 0, len(row))

		for _, button := range row {
			switch {
			case button.Callback != "":
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Callback))
			case button.URL != "":
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonURL(button.Text, button.URL))
			default:
				// Stub for a button with neither callback nor url
				// It should be a non-action button
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, "#"))
			}
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
