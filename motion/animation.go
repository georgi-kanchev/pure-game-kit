/*
Provides a solution for smooth number transitions in the form of the universally known Tweens
(an animation term, short for "inbetween"), that make use of custom curves or pre-made easings,
describing the movement. Also provides a solution for animation sequences - a collection of any data,
iterated over time.
*/
package motion

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

type Animation[T any] struct {
	Items               []T
	ItemsPerSecond      float32
	IsLooping, IsPaused bool

	Time float32
}

func NewAnimation[T any](itemsPerSecond float32, loop bool, items ...T) Animation[T] {
	return Animation[T]{Items: items, ItemsPerSecond: itemsPerSecond, IsLooping: loop}
}

//=================================================================

func (a *Animation[T]) Update() {
	if !a.IsPaused {
		a.Time += internal.DeltaTime
	}

	var duration = a.Duration()
	if a.Time >= duration {
		a.Time = condition.If(a.IsLooping, 0, duration)
	}
}

func (a *Animation[T]) SetDuration(seconds float32) {
	a.ItemsPerSecond = float32(len(a.Items)) / seconds
}
func (a *Animation[T]) SetIndex(index int) {
	index = number.Limit(index, 0, len(a.Items)-1)
	a.Time = number.Map(float32(index), 0, float32(len(a.Items)), 0, a.Duration())
}

//=================================================================

func (a *Animation[T]) Item() *T {
	return &a.Items[a.Index()]
}
func (a *Animation[T]) Index() int {
	return int(number.Map(a.Time, 0, a.Duration(), 0, float32(len(a.Items))))
}
func (a *Animation[T]) Duration() float32 {
	return float32(len(a.Items)) / a.ItemsPerSecond
}

func (a *Animation[T]) IsFinished() bool {
	return a.Time == a.Duration()
}
func (a *Animation[T]) IsPlaying() bool {
	return !a.IsFinished() && !a.IsPaused
}
