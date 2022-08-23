package spec

import (
	"io"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Spec is a base struct for a bot specification file.
type Spec struct {
	Bot *Bot `yaml:"bot"`
}

// Bot spec includes bot configuration and handlers.
type Bot struct {
	Token    string            `yaml:"token" env:"BOT_TOKEN"`
	State    map[string]string `yaml:"state"`
	Handlers []*Handler        `yaml:"handlers"`
}

// Handler specification declares bot handlers.
type Handler struct {
	Trigger *Trigger `yaml:"on"`
	Replies []*Reply `yaml:"reply"`
	State   *State   `yaml:"state"`
}

var ErrNoTriggerConfig = errors.New("no trigger configuration")

func (h *Handler) validate() (errs []error) {
	errs = make([]error, 0)
	if h.Trigger == nil {
		errs = append(errs, ErrNoTriggerConfig)
	} else {
		errs = append(errs, h.Trigger.validate()...)
	}
	for _, reply := range h.Replies {
		errs = append(errs, reply.validate()...)
	}
	return
}

func (s *Spec) parseYAML(dec *yaml.Decoder) error {
	return dec.Decode(s)
}

func (s *Spec) parseEnv() error {
	return env.Parse(s)
}

func (s *Spec) validate() (errs []error) {
	errs = make([]error, 0)
	if s.Bot == nil {
		errs = append(errs, ErrInvalidSpec)
		return
	}
	if s.Bot.Token == "" {
		errs = append(errs, ErrNoTokenProvided)
	}
	if len(s.Bot.Handlers) == 0 {
		errs = append(errs, ErrNoHandlersConfig)
	}
	for _, handler := range s.Bot.Handlers {
		errs = append(errs, handler.validate()...)
	}
	return
}

// ParseYaml decodes YAML input into a Spec struct.
func ParseYaml(r io.Reader) (*Spec, error) {
	var spec Spec
	if err := spec.parseYAML(yaml.NewDecoder(r)); err != nil {
		return nil, errors.Wrap(err, "parse yaml")
	}
	if err := spec.parseEnv(); err != nil {
		return nil, errors.Wrap(err, "parse env")
	}
	if spec.Bot.State == nil {
		spec.Bot.State = make(map[string]string)
	}
	return &spec, nil
}

// Validate specification.
func (s *Spec) Validate() []error {
	return s.validate()
}
