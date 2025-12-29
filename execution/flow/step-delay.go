package flow

import "pure-game-kit/utility/time"

type StepDelay struct {
	startTime, delay float32
}

func NowWaitForDelay(seconds float32) *StepDelay {
	return &StepDelay{delay: seconds, startTime: -1}
}

func (s *StepDelay) Continue(*Sequence) bool {
	var runtime = time.Runtime()

	if s.startTime < 0 {
		s.startTime = runtime
	}

	if runtime > s.startTime+s.delay {
		s.startTime = -1
		return true
	}
	return false
}
