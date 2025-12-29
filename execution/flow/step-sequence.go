package flow

type StepSequence struct {
	seq *Sequence
}

func NowWaitForSequence(sequence *Sequence) *StepSequence {
	return &StepSequence{seq: sequence}
}

func (s *StepSequence) Continue(*Sequence) bool {
	return s.seq.currIndex < 0 || s.seq.currIndex >= len(s.seq.steps)
}
