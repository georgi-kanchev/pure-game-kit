package flow

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/collection"
)

type StepSignal struct {
	signal string
}

func WaitForSignal(signal string) *StepSignal {
	return &StepSignal{signal: signal}
}

func (step *StepSignal) Continue() bool {
	if collection.Contains(internal.FlowSignals, step.signal) {
		var i = collection.IndexOf(internal.FlowSignals, step.signal)
		internal.FlowSignals = collection.RemoveAt(internal.FlowSignals, i, i+1)
		return true
	}
	return false
}
