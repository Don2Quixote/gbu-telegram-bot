# Go Blog Updates - Telegram Bot
This service consumes events about new posts in go blog ([go.dev](https://go.dev)) from message broker ([rabbitmq](https://www.rabbitmq.com/)) ([gbu-scanner service](https://github.com/don2quixote/gbu-scanner) publishes these events) and sends notifications to telegram [bot's](https://core.telegram.org/bots/api) subscribers.
It uses [PostgreSQL](https://www.postgresql.org/) as a storage for bot's users.

### Consumers:
 - [gbu-telegram-bot](https://github.com/don2quixote/gbu-telegram-bot) service