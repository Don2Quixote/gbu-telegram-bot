package messages

import "time"

// Config contains settings for Messages.
type Config struct {
	// Webhook is configuration for setting up listening for webhook
	// If nil then long polling will be used for getting updates
	Webhook *Webhook
	// SendingDelay is time between outgoing messages to
	// to avoid hitting limits (https://core.telegram.org/bots/faq#broadcasting-to-users)
	SendingDelay time.Duration
}

// Webhook is configuration for setting up webhook.
type Webhook struct {
	// Port is port for launching server to listen for requests
	Port int
	// Path is /path where handler will be registered
	Path string
}
