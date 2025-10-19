package flow

import "pure-game-kit/utility/collection"

type StepSignal struct {
	signal string
}

func NowWaitForSignal(signal string) *StepSignal {
	return &StepSignal{signal: signal}
}

func (step *StepSignal) Continue(sequence *Sequence) bool {
	if collection.Contains(sequence.signals, step.signal) {
		var i = collection.IndexOf(sequence.signals, step.signal)
		sequence.signals = collection.RemoveAt(sequence.signals, i, i+1)
		return true
	}
	return false
}
