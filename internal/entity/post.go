package entity

import "time"

// Post is a short data about post in blog, without content.
// Content is available at the URL
type Post struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Author  string    `json:"author"`
	Summary string    `json:"summary"`
	URL     string    `json:"url"`
}

// PostEvent is struct with post and functions to control event's acknowledgement
type PostEvent struct {
	Post Post
	Ack  func() error
	Nack func(requeue bool) error
}
