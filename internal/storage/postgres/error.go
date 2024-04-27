package postgres

import "errors"

var (
	errUserExists    = errors.New("user is exists")
	errUserNotExists = errors.New("user is not exists")
)
