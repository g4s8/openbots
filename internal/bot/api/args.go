package api

import (
	"errors"
	"fmt"

	"github.com/g4s8/openbots/pkg/api"
)

type Argument interface {
	Get(req api.Request) (string, error)
}

type refArg struct {
	ref string
}

var ErrArgumentNotFound = errors.New("argument not found")

func (r *refArg) Get(req api.Request) (string, error) {
	val, ok := req.Payload[r.ref]
	if !ok {
		api.WrapError(ErrArgumentNotFound, api.InvalidRequestDataError,
			fmt.Sprintf("argument `%s` was not present in request payload", r.ref))
	}
	return val, nil
}

func NewRefArg(ref string) Argument {
	return &refArg{ref: ref}
}

type constArg struct {
	val string
}

func (c *constArg) Get(req api.Request) (string, error) {
	return c.val, nil
}

func NewConstArg(val string) Argument {
	return &constArg{val: val}
}
