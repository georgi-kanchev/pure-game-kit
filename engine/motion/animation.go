package motion

import (
	"math"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/seconds"
)

type Animation[T any] struct {
	Items               []T
	ItemsPerSecond      float32
	IsLooping, IsPaused bool

	startTime, lastUpdateTime float32
}

func NewAnimation[T any](itemsPerSecond float32, loop bool, items ...T) Animation[T] {
	return Animation[T]{
		Items: items, ItemsPerSecond: itemsPerSecond, IsLooping: loop, startTime: seconds.GetRuntime()}
}

func (sequence *Animation[T]) Duration() float32 {
	var count = float32(len(sequence.Items))
	return count / sequence.ItemsPerSecond
}

func (sequence *Animation[T]) SetDuration(seconds float32) {
	sequence.ItemsPerSecond = float32(len(sequence.Items)) / seconds
}
func (sequence *Animation[T]) SetIndex(fromIndex int) {
	fromIndex = number.LimitInt(fromIndex, 0, len(sequence.Items)-1)
	var newTime = float32(fromIndex) / sequence.ItemsPerSecond
	sequence.startTime = seconds.GetRuntime() - newTime
}
func (sequence *Animation[T]) SetTime(fromTime float32) {
	sequence.startTime = seconds.GetRuntime() - fromTime
}

func (sequence *Animation[T]) CurrentItem() *T {
	return &sequence.Items[sequence.CurrentIndex()]
}
func (sequence *Animation[T]) CurrentIndex() int {
	var progress = float64(sequence.update())
	var count = float64(len(sequence.Items))
	return int(math.Min(progress*count, count-1))
}
func (sequence *Animation[T]) CurrentTime() float32 {
	var progress = sequence.update()
	var count = float64(len(sequence.Items))
	return progress * float32(count) / sequence.ItemsPerSecond
}

// #region private

func (sequence *Animation[T]) update() float32 {
	var runtime = seconds.GetRuntime()
	var progress = (runtime - sequence.startTime) / sequence.Duration()

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

// #endregion
