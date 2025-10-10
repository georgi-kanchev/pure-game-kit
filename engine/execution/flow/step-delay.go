package flow

import "pure-kit/engine/utility/time"

type StepDelay struct {
	startTime, delay float32
}

func NowWaitForDelay(seconds float32) *StepDelay {
	return &StepDelay{delay: seconds, startTime: -1}
}

func (step *StepDelay) Continue(*Sequence) bool {
	var runtime = time.Runtime()

	if step.startTime < 0 {
		step.startTime = runtime
	}

	if runtime > step.startTime+step.delay {
		step.startTime = -1
		return true
	}
	return false
}
