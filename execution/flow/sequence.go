/*
Mimics the async/await paradigm by providing predetermined steps (rules) and navigating them.
Makes it easier to delay code execution linearly (without nesting) by a known period of time or by
waiting for a specific signal. May also be used as an advanced state machine that executes states
according to different rules or logic, instead of constantly pumping updates.
*/
package flow

import "pure-game-kit/internal"

type Sequence struct {
	steps                []Step
	currIndex, prevIndex int
	stepStartedAt        float32
	signals              []string
}

type Step interface{ Continue(*Sequence) bool }

func NewSequence() *Sequence               { return &Sequence{} }
func (s *Sequence) SetSteps(steps ...Step) { s.steps = steps }

//=================================================================

func (s *Sequence) Signal(signal string)      { s.signals = append(s.signals, signal) }
func (s *Sequence) GoToStep(step int)         { s.currIndex = step }
func (s *Sequence) GoToNextStep()             { s.currIndex++ }
func (s *Sequence) IsRunning() bool           { return s.currIndex >= 0 && s.currIndex < len(s.steps) }
func (s *Sequence) CurrentStep() int          { return s.currIndex }
func (s *Sequence) CurrentStepTimer() float32 { return internal.Runtime - s.stepStartedAt }
func (s *Sequence) Update() {
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
