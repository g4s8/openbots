package spec

type API struct {
	Handlers []*ApiHandler `yaml:"handlers"`
}

type ApiHandler struct {
	ID      string       `yaml:"id"`
	Actions []*ApiAction `yaml:"actions"`
}

type ApiAction struct {
	SendMessage *MessageReply `yaml:"send-message"`
	ChatID      OptUint64     `yaml:"chat-id"`
}
