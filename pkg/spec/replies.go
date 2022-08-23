package spec

type Reply struct {
	Message  *MessageReply  `yaml:"message"`
	Callback *CallbackReply `yaml:"callback"`
}

type MessageReply struct {
	Text   string       `yaml:"text"`
	Markup *ReplyMarkup `yaml:"markup"`
}

type InlineButton struct {
	Text     string `yaml:"text"`
	URL      string `yaml:"url"`
	Callback string `yaml:"callback"`
}

type ReplyMarkup struct {
	Keyboard       [][]string       `yaml:"keyboard"`
	InlineKeyboard [][]InlineButton `yaml:"inlineKeyboard"`
}

type CallbackReply struct {
	Text  string `yaml:"text"`
	Alert bool   `yaml:"alert"`
}
