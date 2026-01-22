package service_errors

import "errors"

var (
	ErrEmptyText              = errors.New("text is empty")
	ErrNoWords                = errors.New("no words")
	ErrSessionNotFound        = errors.New("session not found")
	ErrSessionAlreadyFinished = errors.New("session already finished")
)
