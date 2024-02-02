package filters

import (
	"context"

	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ types.EventFilter = (*FilterChain)(nil)

type FilterChain []types.EventFilter

func (c FilterChain) Check(ctx context.Context, upd *telegram.Update) (bool, error) {
	for _, f := range c {
		ok, err := f.Check(ctx, upd)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

type nopFilter struct{}

func (nopFilter) Check(context.Context, *telegram.Update) (bool, error) { return true, nil }

// Fallback filter always returns true
var Fallback = nopFilter{}

func Join(head types.EventFilter, tail ...types.EventFilter) FilterChain {
	if head == nil {
		head = nopFilter{}
	}
	ch := make(FilterChain, len(tail)+1)
	ch[0] = head
	copy(ch[1:], tail)
	return ch
}
