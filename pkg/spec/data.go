package spec

import "errors"

type Data struct {
	Fetch *DataFetch `yaml:"fetch"`
}

func (c *Data) validate() []error {
	if c.Fetch != nil {
		return c.Fetch.validate()
	}
	return nil
}

type DataFetch struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}

func (c *DataFetch) validate() []error {
	if c.Method == "" {
		c.Method = "GET"
	}
	if c.URL == "" {
		return []error{errors.New("data fetch url is required")}
	}
	return nil
}
