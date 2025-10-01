package flow

import (
	"pure-kit/engine/internal"
)

type StepFlow struct {
	name string
}

func NowWaitForFlow(name string) *StepFlow {
	return &StepFlow{name: name}
}

func (step *StepFlow) Continue() bool {
	var seq, has = internal.Flows[step.name]
	if has {
		return seq.CurrIndex < 0 || seq.CurrIndex >= len(seq.Steps)
	}
	return false
}
