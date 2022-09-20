package spec

import (
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Webhook struct {
	URL    *url.URL `yaml:"url"`
	Method string   `yaml:"method"`
	Body   []byte   `yaml:"body"`
}

var ErrWebhookInvalidURL = errors.New("invalid URL")

func (ch *Webhook) UnmarshalYAML(node *yaml.Node) error {
	var internal struct {
		URL    string `yaml:"url"`
		Method string `yaml:"method"`
		Body   string `yaml:"body"`
	}
	if err := node.Decode(&internal); err != nil {
		return errors.Wrap(err, "decode YAML")
	}
	if internal.URL == "" {
		return ErrWebhookInvalidURL
	}
	u, err := url.Parse(internal.URL)
	if err != nil {
		return errors.Wrap(err, "parse URL")
	}
	ch.URL = u
	ch.Method = internal.Method
	if ch.Method == "" {
		ch.Method = "GET"
	}
	if internal.Body != "" {
		ch.Body = []byte(internal.Body)
	}
	return nil
}
