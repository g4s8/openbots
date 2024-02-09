package bot

import (
	"net/http"

	botctx "github.com/g4s8/openbots/internal/bot/ctx"
	"github.com/g4s8/openbots/pkg/types"
	"github.com/rs/zerolog"
)

// Option is a functional option for Bot.
type Option func(*Bot)

// WithHTTPClient option sets HTTP client for bot.
func WithHTTPClient(cli *http.Client) Option {
	return func(b *Bot) {
		b.httpCli = cli
	}
}

// WithAPIAddr option sets API address for bot.
func WithAPIAddr(addr string) Option {
	return func(b *Bot) {
		b.apiAddr = addr
	}
}

// WithLogger option sets logger for bot.
func WithLogger(log zerolog.Logger) Option {
	return func(b *Bot) {
		b.log = log
	}
}

// WithContextProvider option sets context provider for bot.
func WithContextProvider(cp types.ContextProvider) Option {
	return func(b *Bot) {
		b.cp = botctx.NewProvider(cp)
	}
}

// WithStateProvider option sets state provider for bot.
func WithStateProvider(sp types.StateProvider) Option {
	return func(b *Bot) {
		b.state = sp
	}
}

// WithAssets option sets assets provider for bot.
func WithAssets(assets types.Assets) Option {
	return func(b *Bot) {
		b.assets = assets
	}
}

// WithPaymentProviders option sets payment providers for bot.
func WithPaymentProviders(payments types.PaymentProviders) Option {
	return func(b *Bot) {
		b.payments = payments
	}
}

// WithSecrets option sets secrets provider for bot.
func WithSecrets(secrets types.Secrets) Option {
	return func(b *Bot) {
		b.secrets = secrets
	}
}
