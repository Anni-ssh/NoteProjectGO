package serverErr

import "errors"

var (
	ErrDataNil = errors.New("data not found")
)
