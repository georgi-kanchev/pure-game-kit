package flow

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
)

type Sequence struct {
	steps                []Step
	currIndex, prevIndex int
	stepStartedAt        float32
	signals              []string
	hasPump              bool
}

type Step interface{ Continue(*Sequence) bool }

func NewSequence() *Sequence {
	return &Sequence{}
}
func (sequence *Sequence) SetSteps(runInstantly bool, steps ...Step) {
	sequence.currIndex = -1
	sequence.steps = steps

	if runInstantly {
		sequence.Run()
	}
}

//=================================================================

func (sequence *Sequence) Run() {
	sequence.GoToStep(0)

	if !sequence.hasPump {
		sequence.hasPump = true
		condition.CallFor(number.ValueMaximum[float32](), sequence.update)
	}
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
