package interfaces

// Logger abstracts the logging of messages.
type Logger interface {
	Log(message string) error
}
