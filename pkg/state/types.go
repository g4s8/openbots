package state

type stateChange int

const (
	stateChangeDelete stateChange = iota
	stateChangeWrite
)

type changeReport struct {
	added   []string
	deleted []string
}

type reporter interface {
	changes() changeReport
}
