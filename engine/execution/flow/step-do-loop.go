package flow

type StepLoop struct {
	times, curr int
	action      func(i int)
}

func NowDoLoop(times int, action func(i int)) *StepLoop {
	return &StepLoop{times: times, action: action}
}

func (step *StepLoop) Continue() bool {
	if step.curr >= step.times {
		step.curr = 0
		return true
	}

	step.action(step.curr)
	step.curr++
	return false
}
