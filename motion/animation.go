/*
Provides a solution for smooth number transitions in the form of the universally known Tweens
(an animation term, short for "inbetween"), that make use of custom curves or pre-made easings,
describing the movement. Also provides a solution for animation sequences - a collection of any data,
iterated over time.
*/
package motion

import (
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
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

func (a *Animation[T]) SetDuration(seconds float32) {
	a.ItemsPerSecond = float32(len(a.Items)) / seconds
}
func (a *Animation[T]) SetIndex(index int) {
	index = number.Limit(index, 0, len(a.Items)-1)
	var newTime = float32(index) / a.ItemsPerSecond
	a.startTime = time.Runtime() - newTime
}
func (a *Animation[T]) SetTime(seconds float32) {
	a.startTime = runtime() - seconds
}

//=================================================================

func (a *Animation[T]) CurrentItem() *T {
	return &a.Items[a.CurrentIndex()]
}
func (a *Animation[T]) CurrentIndex() int {
	var progress = a.update()
	var count = float32(len(a.Items))
	return int(number.Smallest(progress*count, count-1))
}
func (a *Animation[T]) CurrentTime() float32 {
	var progress = a.update()
	var count = float64(len(a.Items))
	return progress * float32(count) / a.ItemsPerSecond
}
func (a *Animation[T]) CurrentDuration() float32 {
	var count = float32(len(a.Items))
	return count / a.ItemsPerSecond
}

func (a *Animation[T]) IsFinished() bool {
	var progress = a.update()
	return progress == 1
}
func (a *Animation[T]) IsPlaying() bool {
	return !a.IsFinished() && !a.IsPaused
}

//=================================================================
// private

func (a *Animation[T]) update() float32 {
	var runtime = time.Runtime()
	var progress = (runtime - a.startTime) / a.CurrentDuration()

	if a.IsPaused && runtime != a.lastUpdateTime {
		a.startTime += runtime - a.lastUpdateTime
	}

	if progress >= 1 {
		if a.IsLooping {
			a.startTime = runtime
			progress = 0
		} else {
			progress = 1
		}
	}

	a.lastUpdateTime = runtime
	return progress
}
func runtime() float32 { return time.Runtime() } // freeing up "seconds" as var/param name
