package errSql

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNoteNotFound = errors.New("note not found")
)
