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

// ParseYaml decodes YAML input into a Spec struct.
func ParseYaml(r io.Reader) (*Spec, error) {
	var spec Spec
	if err := yaml.NewDecoder(r).Decode(&spec); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	if err := env.Parse(spec.Bot); err != nil {
		return nil, errors.Wrap(err, "env")
	}
	if spec.Bot.State == nil {
		spec.Bot.State = make(map[string]string)
	}
	return &spec, nil
}

// Validate specification.
func (s *Spec) Validate() error {
	if s.Bot == nil {
		return ErrInvalidSpec
	}
	if s.Bot.Token == "" {
		return ErrNoTokenProvided
	}
	if len(s.Bot.Handlers) == 0 {
		return ErrNoHandlersConfig
	}
	return nil
}
