package adaptors

import (
	"context"

	"github.com/g4s8/openbots/internal/bot/handlers"
	"github.com/g4s8/openbots/pkg/spec"
	"github.com/g4s8/openbots/pkg/types"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.uber.org/multierr"
)

func MessageRepply(bot *telegram.BotAPI,
	sp types.StateProvider, secrets types.Secrets, s *spec.MessageReply,
	log zerolog.Logger,
) (*handlers.MessageReply, error) {
	var modifiers []handlers.MessageModifier
	if s.Markup != nil && len(s.Markup.Keyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithKeyboard(s.Markup.Keyboard))
	}
	if s.Markup != nil && len(s.Markup.InlineKeyboard) > 0 {
		modifiers = append(modifiers, handlers.MessageWithInlineKeyboard(
			inlineKeyboardFromSpec(s.Markup.InlineKeyboard)))
	}
	if s.ParseMode != "" {
		modifiers = append(modifiers, handlers.MessageWithParseMode(string(s.ParseMode)))
	}
	var templater handlers.Templater
	switch s.Template {
	case spec.TemplateDefault:
		templater = handlers.NewDefaultTemplate
	case spec.TemplateGo:
		templater = handlers.NewGoTemplate
	case spec.TemplateNo:
		templater = handlers.NewNoTemplate
	default:
		templater = handlers.NewDefaultTemplate
	}
	tpl, err := templater(s.Text)
	if err != nil {
		return nil, errors.Wrap(err, "create template")
	}
	return handlers.NewMessageReply(bot, sp, secrets, tpl, log, modifiers...), nil
}

func CallbackReply(sp types.StateProvider, secrets types.Secrets, s *spec.CallbackReply) *handlers.CallbackReply {
	return handlers.NewCallbackReply(sp, secrets, s.Text, s.Alert)
}

func inlineKeyboardFromSpec(bts [][]spec.InlineButton) (res handlers.InlineKeyboard) {
	res = make(handlers.InlineKeyboard, len(bts))
	for i, row := range bts {
		res[i] = make([]handlers.InlineButton, len(row))
		for j, btn := range row {
			res[i][j].Text = btn.Text
			res[i][j].URL = btn.URL
			res[i][j].Callback = btn.Callback
		}
	}
	return
}

type multiHandler struct {
	handlers []types.Handler
}

func (h *multiHandler) Handle(ctx context.Context, upd *telegram.Update, bot *telegram.BotAPI) error {
	var merr error
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, upd, bot); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	if merr != nil {
		return errors.Wrap(merr, "one or more handler error")
	}
	return nil
}

func Replies(bot *telegram.BotAPI, sp types.StateProvider, secrets types.Secrets, assets types.Assets, payments types.PaymentProviders,
	r []*spec.Reply, log zerolog.Logger,
) (types.Handler, error) {
	var handlers []types.Handler
	for _, reply := range r {
		if reply.Message != nil {
			h, err := MessageRepply(bot, sp, secrets, reply.Message, log)
			if err != nil {
				return nil, errors.Wrap(err, "create message reply handler")
			}
			handlers = append(handlers, h)
		}
		if reply.Callback != nil {
			handlers = append(handlers, CallbackReply(sp, secrets, reply.Callback))
		}
		if reply.Edit != nil {
			h, err := newEdit(reply.Edit, sp, secrets, log)
			if err != nil {
				return nil, errors.Wrap(err, "create edit handler")
			}
			handlers = append(handlers, h)
		}
		if reply.Delete {
			handlers = append(handlers, newDelete(log))
		}
		if reply.Image != nil {
			handlers = append(handlers, newImageReply(reply.Image, assets, log))
		}
		if reply.Document != nil {
			handlers = append(handlers, newDocumentReply(reply.Document, assets, log))
		}
		if reply.Invoice != nil {
			handlers = append(handlers, newInvoice(reply.Invoice, payments, sp, secrets, log))
		}
		if reply.PreCheckout != nil {
			handlers = append(handlers, newPreCheckoutAnswer(reply.PreCheckout, log))
		}
	}
	return &multiHandler{handlers}, nil
}

func Webhook(s *spec.Webhook, sp types.StateProvider, secrets types.Secrets, log zerolog.Logger) types.Handler {
	return handlers.NewWebhook(s.URL, s.Method, s.Headers, s.Data, sp, secrets, log)
}

func newEdit(s *spec.Edit, sp types.StateProvider, secrets types.Secrets, log zerolog.Logger) (types.Handler, error) {
	if s.Message == nil {
		log.Fatal().Msg("Invalid edit spec: message is empty")
	}
	msg := s.Message

	var templater handlers.Templater
	switch msg.Template {
	case spec.TemplateDefault:
		templater = handlers.NewDefaultTemplate
	case spec.TemplateGo:
		templater = handlers.NewGoTemplate
	case spec.TemplateNo:
		templater = handlers.NewNoTemplate
	default:
		templater = handlers.NewDefaultTemplate
	}
	tpl, err := templater(msg.Text)
	if err != nil {
		return nil, errors.Wrap(err, "create template")
	}
	return handlers.NewMessageEdit(msg.Caption, tpl, inlineKeyboardFromSpec(msg.InlineKeyboard),
		sp, secrets, log), nil
}

func newDelete(logger zerolog.Logger) types.Handler {
	return handlers.NewMessageDelete(logger)
}

func newImageReply(s *spec.FileReply, assets types.Assets,
	log zerolog.Logger,
) types.Handler {
	if s.Key != "" {
		return handlers.NewReplyImageFile(s.Key, s.Name, assets, log)
	}
	return nil
}

func newDocumentReply(s *spec.FileReply, assets types.Assets,
	log zerolog.Logger,
) types.Handler {
	if s.Key != "" {
		return handlers.NewReplyDocument(s.Key, s.Name, assets, log)
	}
	return nil
}

func newInvoice(s *spec.Invoice, providers types.PaymentProviders, sp types.StateProvider, secrets types.Secrets,
	log zerolog.Logger,
) types.Handler {
	prices := make([]handlers.InvoicePrice, len(s.Prices))
	for i, p := range s.Prices {
		prices[i] = handlers.InvoicePrice{
			Label:  p.Label,
			Amount: p.Amount,
		}
	}
	return handlers.NewSendInvoice(providers, sp, secrets, s.Provider, handlers.InvoiceConfig{
		Title:       s.Title,
		Description: s.Description,
		Payload:     s.Payload,
		Currency:    s.Currency,
		Prices:      prices,
	}, log.With().Str("handler", "send_invoice").Str("component", "handler").Logger())
}

func newPreCheckoutAnswer(s *spec.PreCheckoutAnswer, log zerolog.Logger) types.Handler {
	return handlers.NewPreCheckout(s.Ok, s.ErrorMessage, log)
}

func Validator(s *spec.Validators, log zerolog.Logger) (types.Handler, error) {
	checks := make([]handlers.Check, len(s.Checks))
	for i, sc := range s.Checks {
		c, err := handlers.ParseCheckString(string(sc))
		if err != nil {
			return nil, errors.Wrap(err, "parse check string")
		}
		checks[i] = c
	}
	return handlers.NewValidator(log.With().Str("handler", "validator").Str("component", "validators").Logger(),
		s.ErrorMessage, checks...), nil
}
