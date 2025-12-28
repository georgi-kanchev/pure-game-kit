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

func (step *StepRepeat) Continue(*Sequence) bool {
	if step.duration > 0 {
		step.currTime += internal.DeltaTime
		step.actionSec()

		if step.currTime > step.duration {
			step.currTime = 0
			return true
		}
		return false
	}

	if step.i >= step.times {
		step.i = 0
		return true
	}

	step.action(step.i)
	step.i++
	return false
}
