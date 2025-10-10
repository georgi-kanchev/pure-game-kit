package flow

type StepSequence struct {
	seq *Sequence
}

func NowWaitForSequence(sequence *Sequence) *StepSequence {
	return &StepSequence{seq: sequence}
}

func (step *StepSequence) Continue(*Sequence) bool {
	return step.seq.currIndex < 0 || step.seq.currIndex >= len(step.seq.steps)
}
