package spec

import "errors"

var (
	ErrInvalidSpec      = errors.New("invalid configuration")
	ErrNoTokenProvided  = errors.New("no token provided")
	ErrNoHandlersConfig = errors.New("no handlers configured")
)
