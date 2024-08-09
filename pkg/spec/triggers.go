package spec

// Trigger is a handler trigger whcich configures when the handler should be
// executed.
type Trigger struct {
	Message      *MessageTrigger
	Callback     *CallbackTrigger
	Context      string
	PreCheckout  *PreCheckoutTrigger
	PostCheckout *PostCheckoutTrigger
	State        []StateCondition
	Fallback     bool
}

// TriggerType is a type of trigger.
//
//go:generate stringer -type=TriggerType
type TriggerType int

// Trigger types.
const (
	TriggerTypeMessage TriggerType = 1 + iota
	TriggerTypeCallback
	TriggerTypeContext
	TriggerTypePreCheckout
	TriggerTypePostCheckout
	TriggerTypeState
	TriggerTypeFallback
)

// Types returns a list of trigger types.
func (t *Trigger) Types() []TriggerType {
	var typ []TriggerType
	if t.Message != nil {
		typ = append(typ, TriggerTypeMessage)
	}
	if t.Callback != nil {
		typ = append(typ, TriggerTypeCallback)
	}
	if t.Context != "" {
		typ = append(typ, TriggerTypeContext)
	}
	if t.PreCheckout != nil {
		typ = append(typ, TriggerTypePreCheckout)
	}
	if t.PostCheckout != nil {
		typ = append(typ, TriggerTypePostCheckout)
	}
	if len(t.State) > 0 {
		typ = append(typ, TriggerTypeState)
	}
	if t.Fallback {
		typ = append(typ, TriggerTypeFallback)
	}
	return typ
}

type MessageTrigger struct {
	Text    []string
	Command string
}

type CallbackTrigger struct {
	Data string
}

type StateCondition struct {
	Key     string  `yaml:"key"`
	Present OptBool `yaml:"present"`
	Eq      string  `yaml:"eq"`
	NEq     string  `yaml:"neq"`
}
