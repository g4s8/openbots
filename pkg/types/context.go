package types

type Context string

var EmptyContext = Context("")

func (c *Context) Set(value string) {
	*c = Context(value)
}

func (c *Context) Delete(value string) {
	if *c == Context(value) {
		*c = ""
	}
}

func (c Context) Check(value string) bool {
	return string(c) == value
}
