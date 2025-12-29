package motion

import (
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
)

type Tween struct {
	IsPaused                  bool
	tweens                    []action
	currIndex                 int
	startTime, lastUpdateTime float32
}

func NewTween(startingItems ...float32) *Tween {
	if len(startingItems) == 0 {
		return &Tween{}
	}

	var newChain = &Tween{tweens: make([]action, 0)}
	var itemsCopy = make([]float32, len(startingItems))
	copy(itemsCopy, startingItems)

	var newValue = action{
		duration: 0,
		from:     itemsCopy,
		to:       itemsCopy,
		current:  itemsCopy,
	}
	newChain.tweens = append(newChain.tweens, newValue)
	newChain.Restart()

	return newChain
}

//=================================================================

func (t *Tween) Restart() {
	t.currIndex = 0
	t.startTime = time.Runtime()
	t.IsPaused = false

	for i := range t.tweens {
		var action = &t.tweens[i]
		if i == 0 {
			copy(action.current, action.from)
			continue
		}

		copy(action.current, t.tweens[i-1].to)
	}
}

//=================================================================

func (t *Tween) GoTo(duration float32, easing func(progress float32) float32, targets ...float32) *Tween {
	if len(t.tweens) == 0 {
		return t
	}

	var lastTween = t.last()
	var copyFrom = make([]float32, len(lastTween.to))
	var copyTo = make([]float32, len(targets))
	var copyCurr = make([]float32, len(lastTween.to))
	copy(copyFrom, lastTween.to)
	copy(copyTo, targets)
	copy(copyCurr, lastTween.to)

	var newTween = action{
		duration: duration,
		from:     copyFrom,
		to:       copyTo,
		current:  copyCurr,
		easing:   easing,
	}
	t.tweens = append(t.tweens, newTween)

	return t
}
func (t *Tween) Wait(seconds float32) *Tween {
	if len(t.tweens) > 0 {
		var lastTween = t.last()
		t.GoTo(seconds, nil, lastTween.to...)
	}
	return t
}

func (t *Tween) CurrentValues() []float32 {
	var runtime = time.Runtime()
	var elapsed = runtime - t.startTime

	if t.IsFinished() {
		return t.last().current
	}

	var act = &t.tweens[t.currIndex]
	var actionDone = elapsed > act.duration

	if len(t.tweens) == 0 || len(act.from) != len(act.to) || len(act.from) != len(act.current) {
		return []float32{}
	}
	if t.IsPaused && runtime != t.lastUpdateTime {
		t.startTime += runtime - t.lastUpdateTime
	}
	t.lastUpdateTime = runtime

	if t.IsPaused {
		return act.current
	}

	for i := range act.current {
		if actionDone {
			act.current[i] = act.to[i]
			continue
		}
		var progress = elapsed / act.duration
		var ease float32 = progress

		if act.easing != nil {
			ease = act.easing(progress)
		}

		act.current[i] = number.Map(ease, 0, 1, act.from[i], act.to[i])
	}

	if actionDone {
		t.startTime = runtime
		t.currIndex++
	}

	return act.current
}

func (t *Tween) IsFinished() bool {
	return t.currIndex >= len(t.tweens)
}
func (t *Tween) IsPlaying() bool {
	return !t.IsFinished() && !t.IsPaused
}

// =================================================================
// private

type action struct {
	duration          float32
	from, to, current []float32
	easing            func(progress float32) float32
}

func (t *Tween) last() *action { return &t.tweens[len(t.tweens)-1] }
