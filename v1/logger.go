package v1

import "log"

// Logger ...
type Logger interface {
	Error(err error)
}

type logger struct {}

func (l *logger) Error(err error) {
	log.Print(err)
}

// NewLogger ...
func NewLogger() Logger {
	return &logger{}
}