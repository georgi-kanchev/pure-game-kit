package flow

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

type StepRepeat struct {
	times, i int
	action   func(i int)

	duration, currTime float32
	actionSec          func()
}

func NowDoAndRepeat(times int, action func(i int)) *StepRepeat {
	return &StepRepeat{times: times, action: action}
}
func NowDoAndKeepRepeating(action func()) *StepRepeat {
	return &StepRepeat{duration: number.ValueMaximum[float32](), actionSec: action}
}
func NowDoAndRepeatFor(seconds float32, action func()) *StepRepeat {
	return &StepRepeat{duration: seconds, actionSec: action}
}

func (s *StepRepeat) Continue(*Sequence) bool {
	if s.duration > 0 {
		s.currTime += internal.DeltaTime
		s.actionSec()

		if s.currTime > s.duration {
			s.currTime = 0
			return true
		}
		return false
	}

	if s.i >= s.times {
		s.i = 0
		return true
	}

	s.action(s.i)
	s.i++
	return false
}
