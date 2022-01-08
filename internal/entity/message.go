package entity

import "strconv"

// IncomingMessage is struct that represents incoming to bot command
type IncomingMessage struct {
	// From is struct with information about sender
	From MessageSender
	// Text is message's text
	Text string
}

// MessageSender is struct that represents message's sender
type MessageSender struct {
	// ID is sender's id
	ID int64
	// Username is sender's username
	Username string
	// Name is sender's name
	Name string
}

// String method to implement fmt.Stringer interface
func (s MessageSender) String() string {
	if s.Username != "" {
		return "@" + s.Username
	}
	base := 10 // Konche-linter (gomnd)
	return strconv.FormatInt(s.ID, base)
}

// OutgoingMessage is struct that represents outgoing telegram message
// Keyboard and InlineKeyboard are not required, but can't be set both in one OutgoingMessage
// If both Keyboard and InlineKeyboard are not nil, only Keyboard is used
type OutgoingMessage struct {
	Text           string
	Keyboard       *Keyboard
	InlineKeyboard *InlineKeyboard
}

// Keyboard is struct that represents telegram chat's keyboard
// Keyboard[row][column] = button's text
type Keyboard [][]string

// InlineKeyboard is struct that represents telegram chat's keyboard attached to a message
// InlineKeyboard[row][column] = inline button with either URL or callback
type InlineKeyboard [][]InlineButton

// InlineButton is button in InlineKeyboard
// Either URL or callback should be specified. If both are not empty strings, Callback used
type InlineButton struct {
	Text     string
	Callback string
	URL      string
}
