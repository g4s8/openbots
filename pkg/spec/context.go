package spec

import "errors"

type Context struct {
	Set    string `yaml:"set"`
	Delete string `yaml:"delete"`
}

func (c *Context) validate() []error {
	if c.Set == "" && c.Delete == "" {
		return []error{errors.New("empty context")}
	}
	return []error{}
}
