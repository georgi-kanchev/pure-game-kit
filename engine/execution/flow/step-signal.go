package flow

import (
	"pure-kit/engine/internal"
	"slices"
)

type StepSignal struct {
	signal string
}

func WaitForSignal(signal string) *StepSignal {
	return &StepSignal{signal: signal}
}

func (step *StepSignal) Continue() bool {
	if slices.Contains(internal.FlowSignals, step.signal) {
		var i = slices.Index(internal.FlowSignals, step.signal)
		internal.FlowSignals = slices.Delete(internal.FlowSignals, i, i+1)
		return true
	}
	return false
}
