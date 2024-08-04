package errs

import (
	"errors"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotExists = errors.New("user does not exist")
	ErrNoteNotExists = errors.New("note does not exist")
)

func CheckCustomErr(err error) bool {
	if errors.Is(err, ErrUserExists) {
		return true
	}
	if errors.Is(err, ErrUserNotExists) {
		return true
	}
	if errors.Is(err, ErrNoteNotExists) {
		return true
	}
	return false
}
