package note

import "errors"

var (
	ErrNotFound      = errors.New("note not found")
	ErrInvalidTitle  = errors.New("title is required")
	ErrInvalidInput  = errors.New("invalid input")
)
