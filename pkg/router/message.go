package router

// Message is an interface to interface with various chat platform
type Message interface {
	// Send sends a message to an ID (user, channel, etc)
	Send(text string, ID string) error

	// Get returns the current message we're processing
	// it should be reflected to whatever type you expect. You should
	// only use this when you have to
	Get() interface{}

	// GetID returns the current chat id, use with Send()
	GetID() string
}
