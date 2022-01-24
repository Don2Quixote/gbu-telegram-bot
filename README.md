# Go Blog Updates - Telegram Bot
This service consumes events about new posts in go blog ([go.dev](https://go.dev)) from message broker ([rabbitmq](https://www.rabbitmq.com/)) ([gbu-scanner service](https://github.com/don2quixote/gbu-scanner) publishes these events) and sends notifications to telegram [bot's](https://core.telegram.org/bots/api) subscribers.
It uses [PostgreSQL](https://www.postgresql.org/) as a storage for bot's users.

## ENV Configuration:
| name                   | type   | description                                                                        |
| ---------------------- | ------ | ---------------------------------------------------------------------------------- |
| TELEGRAM_BOT_TOKEN     | string | Token to authorize in telegram                                                     |
| POSTGRES_HOST          | string | Database host                                                                      |
| POSTGRES_USER          | string | Database user                                                                      |
| POSTGRES_PASS          | string | Database password                                                                  |
| POSTGRES_DB_NAME       | string | Database name                                                                      |
| RABBIT_HOST            | string | Rabbit host                                                                        |
| RABBIT_USER            | string | Rabbit user                                                                        |
| RABBIT_PASS            | string | Rabbit password                                                                    |
| RABBIT_VHOST           | string | Rabbit vhost                                                                       |
| RABBIT_AMQPS           | bool   | Flag to use amqps protocol instead of amqp                                         |
| RABBIT_RECONNECT_DELAY | int    | Delay (seconds) before attempting to reconnect to rabbit after loosing connection  |
| MESSAGES_SENDING_DELAY | int    | Delay (milliseconds) between sending messages in telegram to avoid hitting [limits](https://core.telegram.org/bots/faq#broadcasting-to-users)|

Env template for sourcing is [deployments/local.env](deployments/local.env)
```
$ source deployments/local.env
```

## Makefile commands:
| name | description                                                                            |
| ---- | -------------------------------------------------------------------------------------- |
| lint | Runs linters                                                                           |
| test | Runs tests, but there are no tests                                                     |
| run  | Sources env variables from [deployments/local.env](deployments/local.env) and runs app |
| stat | Prints stats information about project (packages, files, lines, chars count)           |

Директория [scripts](/scripts) содержит скрипты, которые вызываются командами из [Makefile](Makefile)