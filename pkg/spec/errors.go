package spec

import "errors"

var (
	ErrInvalidSpec      = errors.New("invalid configuration")
	ErrNoHandlersConfig = errors.New("no handlers configured")
)
