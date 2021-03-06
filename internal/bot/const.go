package bot

import "gbu-telegram-bot/internal/entity"

var greetingMessage = entity.OutgoingMessage{
	Text: `đ Welcome to Go Blog Updates bot

đĨ I will notify you once something will be posted in [blog](https://go.dev/blog)

âšī¸ You can unsubscribe from updates any moment pressing button or sending command /unsubscribe`,
	Keyboard: &entity.Keyboard{
		[]string{"â Unsubscribe"},
	},
}

var subscribedMessage = entity.OutgoingMessage{
	Text: `â You have subscribed`,
	Keyboard: &entity.Keyboard{
		[]string{"â Unsubscribe"},
	},
}

var unsubscribedMessage = entity.OutgoingMessage{
	Text: `â You have unsubscribed`,
	Keyboard: &entity.Keyboard{
		[]string{"đĨ Subscribe"},
	},
}

var unknownMessage = entity.OutgoingMessage{
	Text: `â I don't understand`,
}

var errorMessage = entity.OutgoingMessage{
	Text: `â Something went wrong`,
}
