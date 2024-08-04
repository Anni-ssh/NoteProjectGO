package errs

import (
	"errors"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user does not exist")
	ErrNoteNotFound  = errors.New("note not found")
)

func CheckCustomErr(err error) bool {
	if errors.Is(err, ErrUserExists) {
		return true
	}
	if errors.Is(err, ErrUserNotExists) {
		return true
	}
	if errors.Is(err, ErrNoteNotFound) {
		return true
	}
	return false
}
