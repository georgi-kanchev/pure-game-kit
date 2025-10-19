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
// setters

func (tween *Tween) Restart() {
	tween.currIndex = 0
	tween.startTime = time.Runtime()
	tween.IsPaused = false

	for i := range tween.tweens {
		var action = &tween.tweens[i]
		if i == 0 {
			copy(action.current, action.from)
			continue
		}

		copy(action.current, tween.tweens[i-1].to)
	}
}

//=================================================================
// getters

func (chain *Tween) GoTo(duration float32, easing func(progress float32) float32, targets ...float32) *Tween {
	if len(chain.tweens) == 0 {
		return chain
	}

	var lastTween = chain.last()
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
	chain.tweens = append(chain.tweens, newTween)

	return chain
}
func (tween *Tween) Wait(seconds float32) *Tween {
	if len(tween.tweens) > 0 {
		var lastTween = tween.last()
		tween.GoTo(seconds, nil, lastTween.to...)
	}
	return tween
}

func (tween *Tween) CurrentValues() []float32 {
	var runtime = time.Runtime()
	var elapsed = runtime - tween.startTime

	if tween.IsFinished() {
		return tween.last().current
	}

	var act = &tween.tweens[tween.currIndex]
	var actionDone = elapsed > act.duration

	if len(tween.tweens) == 0 || len(act.from) != len(act.to) || len(act.from) != len(act.current) {
		return []float32{}
	}
	if tween.IsPaused && runtime != tween.lastUpdateTime {
		tween.startTime += runtime - tween.lastUpdateTime
	}
	tween.lastUpdateTime = runtime

	if tween.IsPaused {
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
		tween.startTime = runtime
		tween.currIndex++
	}

	return act.current
}

func (tween *Tween) IsFinished() bool {
	return tween.currIndex >= len(tween.tweens)
}
func (tween *Tween) IsPlaying() bool {
	return !tween.IsFinished() && !tween.IsPaused
}

// =================================================================
// private

type action struct {
	duration          float32
	from, to, current []float32
	easing            func(progress float32) float32
}

func (tween *Tween) last() *action { return &tween.tweens[len(tween.tweens)-1] }
