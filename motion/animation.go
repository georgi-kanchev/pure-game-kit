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

func (animation *Animation[T]) SetDuration(seconds float32) {
	animation.ItemsPerSecond = float32(len(animation.Items)) / seconds
}
func (animation *Animation[T]) SetIndex(index int) {
	index = number.Limit(index, 0, len(animation.Items)-1)
	var newTime = float32(index) / animation.ItemsPerSecond
	animation.startTime = time.Runtime() - newTime
}
func (animation *Animation[T]) SetTime(seconds float32) {
	animation.startTime = runtime() - seconds
}

//=================================================================

func (animation *Animation[T]) CurrentItem() *T {
	return &animation.Items[animation.CurrentIndex()]
}
func (animation *Animation[T]) CurrentIndex() int {
	var progress = animation.update()
	var count = float32(len(animation.Items))
	return int(number.Smallest(progress*count, count-1))
}
func (animation *Animation[T]) CurrentTime() float32 {
	var progress = animation.update()
	var count = float64(len(animation.Items))
	return progress * float32(count) / animation.ItemsPerSecond
}
func (animation *Animation[T]) CurrentDuration() float32 {
	var count = float32(len(animation.Items))
	return count / animation.ItemsPerSecond
}

func (animation *Animation[T]) IsFinished() bool {
	var progress = animation.update()
	return progress == 1
}
func (animation *Animation[T]) IsPlaying() bool {
	return !animation.IsFinished() && !animation.IsPaused
}

//=================================================================
// private

func (animation *Animation[T]) update() float32 {
	var runtime = time.Runtime()
	var progress = (runtime - animation.startTime) / animation.CurrentDuration()

	if animation.IsPaused && runtime != animation.lastUpdateTime {
		animation.startTime += runtime - animation.lastUpdateTime
	}

	if progress >= 1 {
		if animation.IsLooping {
			animation.startTime = runtime
			progress = 0
		} else {
			progress = 1
		}
	}

	animation.lastUpdateTime = runtime
	return progress
}
func runtime() float32 { return time.Runtime() } // freeing up "seconds" as var/param name
