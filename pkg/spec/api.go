package spec

type API struct {
	Handlers []*ApiHandler `yaml:"handlers"`
}

type ApiHandler struct {
	ID      string       `yaml:"id"`
	Actions []*ApiAction `yaml:"actions"`
}

type ApiAction struct {
	SendMessage *ApiSendMesageAction `yaml:"send-message"`
}

type ApiSendMesageAction struct {
	Text *ApiArg `yaml:"text"`
}

type ApiArg struct {
	Param string `yaml:"param"`
	Value string `yaml:"value"`
}
