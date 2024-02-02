package spec

import (
	"errors"
	"fmt"
	"io"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

// Spec is a base struct for a bot specification file.
type Spec struct {
	Bot *Bot `yaml:"bot"`
}

// Bot spec includes bot configuration and handlers.
//
//go:generate go run github.com/g4s8/envdoc@latest -output ../../env.md
type Bot struct {
	// Token is a Telegram bot token.
	Token    string            `yaml:"token" env:"BOT_TOKEN"`
	Config   *Config           `yaml:"config"`
	State    map[string]string `yaml:"state"`
	Debug    bool              `yaml:"debug"`
	Handlers []*Handler        `yaml:"handlers"`
	Api      *API              `yaml:"api"`
}

// Handler specification declares bot handlers.
type Handler struct {
	Trigger  *Trigger    `yaml:"on"`
	Replies  []*Reply    `yaml:"reply"`
	Webhook  *Webhook    `yaml:"webhook"`
	State    *State      `yaml:"state"`
	Context  *Context    `yaml:"context"`
	Data     *Data       `yaml:"data"`
	Validate *Validators `yaml:"validate"`
}

var ErrNoTriggerConfig = errors.New("no trigger configuration")

func (h *Handler) validate() error {
	var errs []error
	if h.Trigger == nil {
		errs = append(errs, ErrNoTriggerConfig)
	} else {
		errs = append(errs, h.Trigger.validate())
	}
	for _, reply := range h.Replies {
		errs = append(errs, reply.validate()...)
	}
	if h.Context != nil {
		errs = append(errs, h.Context.validate()...)
	}
	if h.Data != nil {
		errs = append(errs, h.Data.validate()...)
	}
	if h.Validate != nil {
		errs = append(errs, h.Validate.validate()...)
	}
	return errors.Join(errs...)
}

func (s *Spec) parseYAML(dec *yaml.Decoder) error {
	return dec.Decode(s)
}

func (s *Spec) parseEnv() error {
	return env.Parse(s)
}

var ErrMultipleFallbacks = errors.New("multiple fallback handlers")

func (s *Spec) validate() error {
	var errs []error
	if s.Bot == nil {
		return ErrInvalidSpec
	}

	if len(s.Bot.Handlers) == 0 {
		errs = append(errs, ErrNoHandlersConfig)
	}
	var hasFallback bool
	for _, h := range s.Bot.Handlers {
		fb := h.Trigger != nil && h.Trigger.Fallback
		if fb && hasFallback {
			errs = append(errs, ErrMultipleFallbacks)
		}
		if h.Trigger != nil && h.Trigger.Fallback {
			hasFallback = true
		}
	}
	for _, handler := range s.Bot.Handlers {
		errs = append(errs, handler.validate())
	}
	// TODO: move from here or rename method
	if s.Bot.Config == nil {
		s.Bot.Config = &Config{
			Persistence: &PersistenceConfig{
				Type: MemoryPersistence,
			},
		}
	}
	return errors.Join(errs...)
}

// ParseYaml decodes YAML input into a Spec struct.
func ParseYaml(r io.Reader) (*Spec, error) {
	var spec Spec
	if err := spec.parseYAML(yaml.NewDecoder(r)); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	if err := spec.parseEnv(); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	if spec.Bot.State == nil {
		spec.Bot.State = make(map[string]string)
	}
	return &spec, nil
}

// Validate specification.
func (s *Spec) Validate() error {
	return s.validate()
}
