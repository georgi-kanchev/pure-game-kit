package flow

import (
	"pure-kit/engine/internal"
)

type StepFlow struct {
	name string
}

func WaitForAnotherFlow(name string) *StepFlow {
	return &StepFlow{name: name}
}

func (step *StepFlow) Continue() bool {
	var seq, has = internal.Flows[step.name]
	if has {
		return seq.CurrentIndex < 0 || seq.CurrentIndex >= len(seq.Steps)
	}
	return false
}
