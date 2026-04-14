// Provides a solution for smooth number transitions in the form of the universally known Tweens
// (an animation term, short for "inbetween"), that make use of custom curves or pre-made easings,
// describing the movement. Also provides a solution for animation sequences that does not hold its data (frames),
// just iterates their count over time.
package motion

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

type Animation struct {
	ItemCount           int
	ItemsPerSecond      float32
	IsLooping, IsPaused bool

	Time float32
}

func NewAnimation(itemCount int, itemsPerSecond float32, loop bool) Animation {
	return Animation{ItemCount: itemCount, ItemsPerSecond: itemsPerSecond, IsLooping: loop}
}

//=================================================================

func (a *Animation) Update() {
	if !a.IsPaused {
		a.Time += internal.DeltaTime
	}

	var duration = a.Duration()
	if a.Time >= duration {
		a.Time = condition.If(a.IsLooping, 0, duration)
	}
}

func (a *Animation) SetDuration(seconds float32) {
	a.ItemsPerSecond = float32(a.ItemCount) / seconds
}
func (a *Animation) SetIndex(index int) {
	index = number.Limit(index, 0, a.ItemCount-1)
	a.Time = number.Map(float32(index), 0, float32(a.ItemCount), 0, a.Duration())
}

//=================================================================

func (a *Animation) Index() int {
	return int(number.Map(a.Time, 0, a.Duration(), 0, float32(a.ItemCount)))
}
func (a *Animation) Duration() float32 {
	return float32(a.ItemCount) / a.ItemsPerSecond
}

func (a *Animation) IsFinished() bool {
	return a.Time == a.Duration()
}
func (a *Animation) IsPlaying() bool {
	return !a.IsFinished() && !a.IsPaused
}
