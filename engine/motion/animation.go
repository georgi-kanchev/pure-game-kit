package motion

import (
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/time"
)

type Animation[T any] struct {
	Items               []T
	ItemsPerSecond      float32
	IsLooping, IsPaused bool

	startTime, lastUpdateTime float32
}

func NewAnimation[T any](itemsPerSecond float32, loop bool, items ...T) Animation[T] {
	return Animation[T]{
		Items: items, ItemsPerSecond: itemsPerSecond, IsLooping: loop, startTime: time.Runtime()}
}

//=================================================================
// setters

func (sequence *Animation[T]) SetDuration(seconds float32) {
	sequence.ItemsPerSecond = float32(len(sequence.Items)) / seconds
}
func (sequence *Animation[T]) SetIndex(index int) {
	index = number.LimitInt(index, 0, len(sequence.Items)-1)
	var newTime = float32(index) / sequence.ItemsPerSecond
	sequence.startTime = time.Runtime() - newTime
}
func (sequence *Animation[T]) SetTime(seconds float32) {
	sequence.startTime = runtime() - seconds
}

//=================================================================
// getters

func (sequence *Animation[T]) CurrentItem() *T {
	return &sequence.Items[sequence.CurrentIndex()]
}
func (sequence *Animation[T]) CurrentIndex() int {
	var progress = sequence.update()
	var count = float32(len(sequence.Items))
	return int(number.Smallest(progress*count, count-1))
}
func (sequence *Animation[T]) CurrentTime() float32 {
	var progress = sequence.update()
	var count = float64(len(sequence.Items))
	return progress * float32(count) / sequence.ItemsPerSecond
}
func (sequence *Animation[T]) CurrentDuration() float32 {
	var count = float32(len(sequence.Items))
	return count / sequence.ItemsPerSecond
}

func (sequence *Animation[T]) IsFinished() bool {
	var progress = sequence.update()
	return progress == 1
}
func (sequence *Animation[T]) IsPlaying() bool {
	return !sequence.IsFinished() && !sequence.IsPaused
}

//=================================================================
// private

func (sequence *Animation[T]) update() float32 {
	var runtime = time.Runtime()
	var progress = (runtime - sequence.startTime) / sequence.CurrentDuration()

	if sequence.IsPaused && runtime != sequence.lastUpdateTime {
		sequence.startTime += runtime - sequence.lastUpdateTime
	}

	if progress >= 1 {
		if sequence.IsLooping {
			sequence.startTime = runtime
			progress = 0
		} else {
			progress = 1
		}
	}

	sequence.lastUpdateTime = runtime
	return progress
}
func runtime() float32 { return time.Runtime() } // freeing up "seconds" as var/param name
