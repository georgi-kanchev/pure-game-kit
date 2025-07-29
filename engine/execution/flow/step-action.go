package flow

type StepDo struct {
	action func()
}

func Do(action func()) *StepDo {
	return &StepDo{action: action}
}

func (step *StepDo) Continue() bool {
	step.action()
	return true
}
