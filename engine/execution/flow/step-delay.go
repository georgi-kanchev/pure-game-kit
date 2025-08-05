package flow

import "pure-kit/engine/utility/seconds"

type StepDelay struct {
	startTime, delay float32
}

func WaitForDelay(seconds float32) *StepDelay {
	return &StepDelay{delay: seconds, startTime: -1}
}

func (step *StepDelay) Continue() bool {
	var runtime = seconds.Runtime()

	if step.startTime < 0 {
		step.startTime = runtime
	}

	if runtime > step.startTime+step.delay {
		step.startTime = -1
		return true
	}
	return false
}
