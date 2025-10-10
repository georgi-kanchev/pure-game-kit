package flow

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/time"
)

type Sequence struct {
	steps                []Step
	currIndex, prevIndex int
	stepStartedAt        float32
	signals              []string
}

type Step interface{ Continue(*Sequence) bool }

func NewSequence(runInstantly bool, steps ...Step) *Sequence {
	var result = &Sequence{steps: steps, currIndex: -1}
	if runInstantly {
		result.Run()
	}
	return result
}

//=================================================================

func (sequence *Sequence) Run() {
	sequence.GoToStep(0)
	condition.CallFor(number.ValueMaximum[float32](), sequence.update)
}
func (sequence *Sequence) Stop() {
	sequence.GoToStep(-1)
}

func (sequence *Sequence) Signal(signal string) {
	sequence.signals = append(sequence.signals, signal)
}
func (sequence *Sequence) GoToStep(step int) {
	sequence.currIndex = step
}
func (sequence *Sequence) GoToNextStep() {
	sequence.currIndex++
}

func (sequence *Sequence) IsRunning() bool {
	return number.IsBetween(sequence.currIndex, 0, len(sequence.steps), true, false)
}

func (sequence *Sequence) CurrentStep() int {
	return sequence.currIndex
}

// useful for time tracking a continuously looping step
func (sequence *Sequence) CurrentStepTimer() float32 {
	return time.Runtime() - sequence.stepStartedAt
}

//=================================================================
// private

func (sequence *Sequence) update(float32) {
	if sequence.prevIndex != sequence.currIndex {
		sequence.stepStartedAt = internal.Runtime
	}
	sequence.prevIndex = sequence.currIndex

	var prev = sequence.currIndex // this checks if we changed index inside the step itself, skip increment if so
	var validIndex = sequence.currIndex >= 0 && sequence.currIndex < len(sequence.steps)
	var keepGoing = validIndex && sequence.steps[sequence.currIndex].Continue(sequence)

	if keepGoing && prev == sequence.currIndex {
		sequence.currIndex++
	}
}
