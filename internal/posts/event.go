package posts

import "gbu-telegram-bot/internal/entity"

// PostEvent is implementation for bot.PostEvent interface.
type PostEvent struct {
	post entity.Post
	ack  func() error
	nack func(requeue bool) error
}

func (e PostEvent) Post() entity.Post       { return e.post }
func (e PostEvent) Ack() error              { return e.ack() }
func (e PostEvent) Nack(requeue bool) error { return e.nack(requeue) }
