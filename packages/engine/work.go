package engine

import "pure-game-kit/packages/internal"

type WorkId byte

func NewWork(work func()) WorkId {
	var id = len(internal.Work) + 1
	internal.Work[byte(id)] = work
	return WorkId(id)
}

func (l WorkId) Start() {
	internal.Working[byte(l)] = true
	internal.WorkQueue = append(internal.WorkQueue, byte(l))
}
func (l WorkId) IsWorking() bool {
	return internal.Working[byte(l)]
}
func (l WorkId) IsJustFinished() bool {
	return internal.WorkJustFinished[byte(l)]
}
