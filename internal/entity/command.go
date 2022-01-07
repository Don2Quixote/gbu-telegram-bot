package entity

// Command is command's action
type Command int8

const (
	// CommandUnknown is invalid action
	CommandUnknown Command = 0
	// CommandStart when user starts dialog with bot
	CommandStart Command = 1
	// CommandSubscribe for subscribing user for notifications
	CommandSubscribe Command = 2
	// CommandUnsubscribe for unsubscribing user from notifications
	CommandUnsubscribe Command = 3
)

// String method to implement fmt.Stringer interface
func (c Command) String() string {
	switch c {
	case CommandStart:
		return "start"
	case CommandSubscribe:
		return "subscribe"
	case CommandUnsubscribe:
		return "unsubscribe"
	default:
		return "unknown"
	}
}
