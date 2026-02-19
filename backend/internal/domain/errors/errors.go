package errors

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrConflict      = errors.New("conflict")
	ErrAIUnavailable = errors.New("OpenAI API key not configured")
	ErrInternal      = errors.New("internal error")
)
