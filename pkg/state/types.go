package state

type stateChange int

const (
	stateChangeDelete stateChange = iota
	stateChangeWrite
)
