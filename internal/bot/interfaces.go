package bot

import (
	"context"

	"gbu-telegram-bot/internal/entity"
)

// Consumer is interface for interacting with message broker
// from where events about new posts in blog come.
type Consumer interface {
	// Consumer returns channel to which new blog's posts will be sent.
	// Returned chan should be closed when context will be closed.
	Consume(ctx context.Context) (<-chan entity.PostEvent, error)
}

// Users is interface for interacting with storage where
// information about bot's users stored.
type Users interface {
	// Add adds new user.
	Add(ctx context.Context, user entity.User) error
	// Get gets user.
	Get(ctx context.Context, id int64) (entity.User, error)
	// Subscribe marks user as subscribed.
	Subscribe(ctx context.Context, id int64) error
	// Unsubscribe marks user as not subscribed.
	Unsubscribe(ctx context.Context, id int64) error
	// GetSubscribedIDs gets ids of users that should get notifications
	// about new posts in blog.
	GetSubscribedIDs(ctx context.Context) ([]int64, error)
}

// Messages is interface for interacting with telegram's messages.
type Messages interface {
	// Consume returns channel to which telegram bot's
	// incoming messages will be sent.
	// Returned channel should be closed when context is closed.
	Consume(ctx context.Context) (<-chan entity.IncomingMessage, error)
	// Send sends a message to user by ID.
	Send(ctx context.Context, id int64, message entity.OutgoingMessage) error
}
