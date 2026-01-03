package flow

import "pure-game-kit/internal"

type StateMachine struct {
	currentState func()
	timer        float32
}

func NewStateMachine() *StateMachine {
	return &StateMachine{}
}

//=================================================================

func (s *StateMachine) GoToState(state func()) {
	s.currentState = state
	s.timer = 0
}

func (s *StateMachine) UpdateCurrentState() {
	if s.currentState != nil {
		s.timer += internal.DeltaTime
		s.currentState()
	}
}

//=================================================================

func (s *StateMachine) StateTimer() float32 {
	return s.timer
}
