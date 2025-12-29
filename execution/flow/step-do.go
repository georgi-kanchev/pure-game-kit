package flow

type StepDo struct {
	action func()
}

func NowDo(action func()) *StepDo {
	return &StepDo{action: action}
}

func (s *StepDo) Continue(*Sequence) bool {
	s.action()
	return true
}
