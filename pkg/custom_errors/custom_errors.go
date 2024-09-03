package custom_errors

import "errors"

var (
	ErrEmptyPassword = errors.New("empty password")
)
