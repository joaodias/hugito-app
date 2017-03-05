package mocks

type Logger struct{}

func (l *Logger) Log(message string) error {
	return nil
}
