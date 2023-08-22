package spec

// State handler spec.
type State struct {
	Set    map[string]string `yaml:"set"`
	Delete Strings           `yaml:"delete"`
	Ops    []StateUpdateOp   `yaml:"ops"`
}

// StateUpdateOpKind is a kind of state update operation.
type StateUpdateOpKind string

const (
	StateUpdateOpKindSet    StateUpdateOpKind = "set"
	StateUpdateOpKindDelete StateUpdateOpKind = "delete"
	StateUpdateOpKindAdd    StateUpdateOpKind = "add"
	StateUpdateOpKindSub    StateUpdateOpKind = "sub"
	StateUpdateOpKindMul    StateUpdateOpKind = "mul"
	StateUpdateOpKindDiv    StateUpdateOpKind = "div"
)

// StateUpdateOp is a state update operation.
type StateUpdateOp struct {
	Kind  StateUpdateOpKind `yaml:"kind"`
	Key   string            `yaml:"key"`
	Value string            `yaml:"value"`
}
