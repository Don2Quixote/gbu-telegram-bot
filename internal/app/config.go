package app

// appConfig is struct for parsing ENV configuration.
type appConfig struct {
	// TelegramBotToken is token to use Telegram's bot API.
	TelegramBotToken string `config:"TELEGRAM_BOT_TOKEN,required"`
	// PostgresHost is host of postgres database.
	PostgresHost string `config:"POSTGRES_HOST,required"`
	// PostgresUser is user for postgres database.
	PostgresUser string `config:"POSTGRES_USER"`
	// PostgresPass is password for postgres database.
	PostgresPass string `config:"POSTGRES_PASS"`
	// PostgresDBName is database's name in postgres.
	PostgresDBName string `config:"POSTGRES_DB_NAME,required"`
	// RabbitHost is host of rabbitmq.
	RabbitHost string `config:"RABBIT_HOST,required"`
	// RabbitUser is user for rabbitmq.
	RabbitUser string `config:"RABBIT_USER"`
	// RabbitPass is password for rabbitmq.
	RabbitPass string `config:"RABBIT_PASS"`
	// RabbitVhost is vhost in rabbitmq to connect.
	RabbitVhost string `config:"RABBIT_VHOST"`
	// RabbitAmqps flag shows should amqps protocol be used instead of amqp or not.
	RabbitAmqps bool `config:"RABBIT_AMQPS"`
	// RabbitReconnectDelay is delay (in seconds) before attempting to reconnect to rabbit after loosing connection.
	RabbitReconnectDelay int `config:"RABBIT_RECONNECT_DELAY,required"`
	// MessagesSendingDelay is delay (milliseconds) between sending messages in telegram
	// to avoid hitting limits (https://core.telegram.org/bots/faq#broadcasting-to-users).
	MessagesSendingDelay int `config:"MESSAGES_SENDING_DELAY"`
}
