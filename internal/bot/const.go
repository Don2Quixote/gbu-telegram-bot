package bot

import "gbu-telegram-bot/internal/entity"

var greetingMessage = entity.OutgoingMessage{
	Text: `👋 Welcome to Go Blog Updates bot

📥 I will notify you once something will be posted in blog

ℹ️ You can unsubscribe from updates any moment pressing button or sending command /unsubscribe`,
	Keyboard: &entity.Keyboard{
		[]string{"❌ Unsubscribe"},
	},
}

var subscribedMessage = entity.OutgoingMessage{
	Text: `✅ You have subscribed`,
	Keyboard: &entity.Keyboard{
		[]string{"❌ Unsubscribe"},
	},
}

var unsubscribedMessage = entity.OutgoingMessage{
	Text: `✅ You have unsubscribed`,
	Keyboard: &entity.Keyboard{
		[]string{"📥 Subscribe"},
	},
}

var unknownMessage = entity.OutgoingMessage{
	Text: `❌ I don't understand`,
}

var errorMessage = entity.OutgoingMessage{
	Text: `❌ Something went wrong`,
}
