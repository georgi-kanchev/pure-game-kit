// Provides a solution for smooth number transitions in the form of the universally known Tweens
// (an animation term, short for "inbetween"), that make use of custom curves or pre-made easings,
// describing the movement. Also provides a solution for animation sequences - a collection of any data (frames),
// iterated over time.
package motion

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type Animation[T any] struct {
	Frames              []T
	FPS                 float32
	IsLooping, IsPaused bool

	Time float32

	lastUpdateFrame uint64
}

func NewAnimation[T any](fps float32, loop bool, frames ...T) Animation[T] {
	return Animation[T]{Frames: frames, FPS: fps, IsLooping: loop}
}
func NewAnimationFromAsset(assetId assets.AnimationsId, name string, fps float32, loop bool) Animation[assets.ImageId] {
	var frameCount = assetId.FrameCount(name)
	var anim = Animation[assets.ImageId]{Frames: make([]assets.ImageId, frameCount), FPS: fps, IsLooping: loop}
	for i := range frameCount {
		anim.Frames[i] = assetId.Frame(name, i)
	}
	return anim
}

//=================================================================

func (a *Animation[T]) SetDuration(seconds float32) {
	a.tryUpdate()
	a.FPS = float32(len(a.Frames)) / seconds
}
func (a *Animation[T]) SetIndex(index int) {
	a.tryUpdate()
	index = number.Limit(index, 0, len(a.Frames)-1)
	a.Time = number.Map(float32(index), 0, float32(len(a.Frames)), 0, a.Duration())
}

//=================================================================

func (a *Animation[T]) Frame() T {
	a.tryUpdate()
	return a.Frames[a.Index()]
}
func (a *Animation[T]) Index() int {
	a.tryUpdate()
	var index = int(number.Map(a.Time, 0, a.Duration(), 0, float32(len(a.Frames))))
	return number.Limit(index, 0, len(a.Frames)-1)
}
func (a *Animation[T]) Duration() float32 {
	a.tryUpdate()
	return float32(len(a.Frames)) / a.FPS
}

func (a *Animation[T]) IsFinished() bool {
	a.tryUpdate()
	return a.Time == a.Duration()
}
func (a *Animation[T]) IsPlaying() bool {
	a.tryUpdate()
	return !a.IsFinished() && !a.IsPaused
}

// private ========================================================

func (a *Animation[T]) tryUpdate() {
	if a.lastUpdateFrame == internal.Frame {
		return
	}

	a.lastUpdateFrame = internal.Frame

	if !a.IsPaused {
		a.Time += internal.FrameDelta
	}

	var duration = a.Duration()
	if a.Time >= duration {
		if a.IsLooping {
			a.Time = 0
		} else {
			a.Time = duration
		}
	}
}
