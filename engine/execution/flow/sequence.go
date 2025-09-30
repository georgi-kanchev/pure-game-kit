package flow

import (
	"pure-kit/engine/internal"
)

type Step = internal.Step

func NewSequence(name string, steps ...Step) {
	internal.Flows[name] = &internal.Sequence{Steps: steps, CurrentIndex: -1}
}

//=================================================================

func Signal(signal string) {
	internal.FlowSignals = append(internal.FlowSignals, signal)
}
func GoToStep(sequenceName string, index int) {
	var seq, has = internal.Flows[sequenceName]
	if has {
		seq.CurrentIndex = index
	}
}
func Start(sequenceName string) {
	GoToStep(sequenceName, 0)
}
func End(sequenceName string) {
	GoToStep(sequenceName, -1)
}
func Exists(sequenceName string) bool {
	var _, has = internal.Flows[sequenceName]
	return has
}
