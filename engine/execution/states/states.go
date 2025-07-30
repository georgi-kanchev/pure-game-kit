package states

import "pure-kit/engine/internal"

func NewMachine(name string, states ...func()) {
	internal.States[name] = &internal.StateMachine{States: states}
}

func GoToState(machineName string, index int) {
	var machine, has = internal.States[machineName]
	if has {
		machine.CurrentIndex = index
	}
}
