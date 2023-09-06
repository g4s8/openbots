package api

import (
	"fmt"

	"github.com/rs/zerolog"
)

type ErrorKind int

const (
	InvalidRequestDataError ErrorKind = iota
	HandlerFailedError
)

type Error struct {
	base error
	msg  string
	kind ErrorKind
}

func WrapError(err error, kind ErrorKind, msg string) *Error {
	return &Error{
		base: err,
		msg:  msg,
		kind: kind,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.msg, e.base)
}

func (e *Error) Unwrap() error {
	return e.base
}

func (e *Error) MarshalZerologObject(evt *zerolog.Event) {
	evt.AnErr("base", e.base)
	evt.Str("msg", e.msg)
	evt.Int("kind", int(e.kind))
}
