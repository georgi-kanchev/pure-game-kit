// Provides a dead-simple state machine with a timer and an update counter per state.
package flow

import "pure-game-kit/internal"

type StateMachine struct {
	currentState func()
	timer        float32
	counter      int
}

func NewStateMachine() *StateMachine {
	return &StateMachine{}
}

//=================================================================

func (s *StateMachine) GoToState(state func()) {
	s.currentState = state
	s.timer = 0
	s.counter = 0
}
func (s *StateMachine) UpdateCurrentState() {
	if s.currentState != nil {
		s.currentState()
		s.timer += internal.DeltaTime
		s.counter++
	}
}

//=================================================================

func (s *StateMachine) StateTimer() float32 {
	return s.timer
}
func (s *StateMachine) StateUpdateCounter() int {
	return s.counter
}
