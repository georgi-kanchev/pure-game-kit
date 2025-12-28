/*
Mimics the async/await paradigm by providing predetermined steps (rules) and navigating them.
Makes it easier to delay code execution linearly (without nesting) by a known period of time or by
waiting for a specific signal. May also be used as an advanced state machine that executes states
according to different rules, instead of constantly pumping updates.
*/
package flow

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
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
func (s *Sequence) SetSteps(runInstantly bool, steps ...Step) {
	s.currIndex = -1
	s.steps = steps

	if runInstantly {
		s.Run()
	}
}

//=================================================================

func (s *Sequence) Run() {
	s.GoToStep(0)

	if !s.hasPump {
		s.hasPump = true
		condition.CallFor(number.ValueMaximum[float32](), s.update)
	}
}
func (s *Sequence) Stop() {
	s.GoToStep(-1)
}

func (s *Sequence) Signal(signal string) {
	s.signals = append(s.signals, signal)
}
func (s *Sequence) GoToStep(step int) {
	s.currIndex = step
}
func (s *Sequence) GoToNextStep() {
	s.currIndex++
}

func (s *Sequence) IsRunning() bool {
	return number.IsBetween(s.currIndex, 0, len(s.steps), true, false)
}

func (s *Sequence) CurrentStep() int {
	return s.currIndex
}

// useful for time tracking a continuously looping step
func (s *Sequence) CurrentStepTimer() float32 {
	return internal.Runtime - s.stepStartedAt
}

//=================================================================
// private

func (s *Sequence) update(float32) {
	if s.prevIndex != s.currIndex {
		s.stepStartedAt = internal.Runtime
	}
	s.prevIndex = s.currIndex

	var prev = s.currIndex // this checks if we changed index inside the step itself, skip increment if so
	var validIndex = s.currIndex >= 0 && s.currIndex < len(s.steps)
	var keepGoing = validIndex && s.steps[s.currIndex].Continue(s)

	if keepGoing && prev == s.currIndex {
		s.currIndex++
	}
}
