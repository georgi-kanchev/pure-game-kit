package flow

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/time"
)

type Step = internal.Step

func NewSequence(name string, runInstantly bool, steps ...Step) {
	internal.Flows[name] = &internal.Sequence{Steps: steps, CurrIndex: -1}
	Run(name)
}

//=================================================================

func Run(sequenceName string) {
	GoToStep(sequenceName, 0)
}
func Stop(sequenceName string) {
	GoToStep(sequenceName, -1)
}

func Signal(signal string) {
	internal.FlowSignals = append(internal.FlowSignals, signal)
}
func GoToStep(sequenceName string, step int) {
	var seq, has = internal.Flows[sequenceName]
	if has {
		seq.CurrIndex = step
	}
}
func GoToNextStep(sequenceName string) {
	var seq, has = internal.Flows[sequenceName]
	if has {
		seq.CurrIndex++
	}
}

func IsRunning(sequenceName string) bool {
	var seq, has = internal.Flows[sequenceName]
	if has && number.IsBetween(seq.CurrIndex, 0, len(seq.Steps), true, false) {
		return true
	}
	return false
}
func IsExisting(sequenceName string) bool {
	var _, has = internal.Flows[sequenceName]
	return has
}

func CurrentStep(sequenceName string) int {
	var seq, has = internal.Flows[sequenceName]
	if !has {
		return -1
	}
	return seq.CurrIndex
}

// useful for time tracking a continuously looping step
func CurrentStepTimer(sequenceName string) float32 {
	var seq, has = internal.Flows[sequenceName]
	if !has {
		return number.NaN()
	}
	return time.Runtime() - seq.StepStartedAt
}
