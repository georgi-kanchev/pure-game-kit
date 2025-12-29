package flow

import "pure-game-kit/utility/collection"

type StepSignal struct {
	signal string
}

func NowWaitForSignal(signal string) *StepSignal {
	return &StepSignal{signal: signal}
}

func (s *StepSignal) Continue(sequence *Sequence) bool {
	if collection.Contains(sequence.signals, s.signal) {
		var i = collection.IndexOf(sequence.signals, s.signal)
		sequence.signals = collection.RemoveAt(sequence.signals, i, i+1)
		return true
	}
	return false
}
