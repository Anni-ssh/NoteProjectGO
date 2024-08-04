package errs

import (
	"errors"
)

var (
	// Entities
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user does not exist")
	ErrNoteNotExists = errors.New("note does not exist")

	// Token
	ErrInvalidTokenType     = errors.New("invalid token type")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)
