package flow

import (
	"pure-kit/engine/internal"
)

type Step = internal.Step

func NewSequence(name string, steps ...Step) {
	internal.Flows[name] = &internal.Sequence{Steps: steps, CurrentIndex: -1}
}

func Signal(signal string) {
	internal.FlowSignals = append(internal.FlowSignals, signal)
}
func GoToStep(name string, index int) {
	var seq, has = internal.Flows[name]
	if has {
		seq.CurrentIndex = index
	}
}
func Start(name string) {
	GoToStep(name, 0)
}
func End(name string) {
	GoToStep(name, -1)
}
