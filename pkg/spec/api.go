package spec

// API specification declares bot handlers.
type API struct {
	// Handlers is a list of API handlers.
	Handlers []*ApiHandler `yaml:"handlers"`
}

// ApiHandler is a handler for API requests.
type ApiHandler struct {
	// ID is a unique handler identifier.
	ID string `yaml:"id"`
	// Actions to perform.
	Actions []*ApiAction `yaml:"actions"`
}

// ApiAction to perform.
type ApiAction struct {
	SendMessage *MessageReply `yaml:"send-message"`
	State       *State        `yaml:"state"`
	Context     *Context      `yaml:"context"`
	// ChatID is a chat identifier for message reply.
	ChatID Uints64 `yaml:"chat-id"`
}
