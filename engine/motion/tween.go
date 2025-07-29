package motion

import (
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/seconds"
)

type Tween struct {
	IsRunning    bool
	tweens       []tween
	currIndex    int
	startRuntime float32
}

func NewTween(startingItems ...float32) *Tween {
	if len(startingItems) == 0 {
		return &Tween{}
	}

	var newChain = &Tween{tweens: make([]tween, 0)}
	var itemsCopy = make([]float32, len(startingItems))
	copy(itemsCopy, startingItems)

	var newValue = tween{
		duration: 0,
		from:     itemsCopy,
		to:       itemsCopy,
		current:  itemsCopy,
	}
	newChain.tweens = append(newChain.tweens, newValue)
	newChain.Restart()

	return newChain
}

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

	var newTween = tween{
		duration: duration,
		from:     copyFrom,
		to:       copyTo,
		current:  copyCurr,
		easing:   easing,
	}
	chain.tweens = append(chain.tweens, newTween)

	return chain
}
func (chain *Tween) Wait(seconds float32) *Tween {
	if len(chain.tweens) > 0 {
		var lastTween = chain.last()
		chain.GoTo(seconds, nil, lastTween.to...)
	}
	return chain
}
func (chain *Tween) Restart() {
	chain.currIndex = 0
	chain.startRuntime = seconds.GetRuntime()
	chain.IsRunning = true

	for i := range chain.tweens {
		var tween = &chain.tweens[i]
		if i == 0 {
			copy(tween.current, tween.from)
			continue
		}

		copy(tween.current, chain.tweens[i-1].to)
	}
}

func (chain *Tween) CurrentValues() []float32 {
	if len(chain.tweens) == 0 {
		return []float32{}
	}

	var tween = &chain.tweens[chain.currIndex]

	if !chain.IsRunning {
		return tween.current
	}

	if len(tween.from) != len(tween.to) || len(tween.from) != len(tween.current) {
		return []float32{}
	}

	var runtime = seconds.GetRuntime()
	var elapsed = runtime - chain.startRuntime
	var tweenDone = elapsed > tween.duration

	for i := range tween.current {
		if tweenDone {
			tween.current[i] = tween.to[i]
			continue
		}
		var progress = elapsed / tween.duration
		var ease float32 = progress

		if tween.easing != nil {
			ease = tween.easing(progress)
		}

		tween.current[i] = number.Map(ease, 0, 1, tween.from[i], tween.to[i])
	}

	if tweenDone {
		chain.startRuntime = runtime
		chain.currIndex++
	}

	var chainDone = chain.currIndex >= int(len(chain.tweens))

	if chainDone {
		chain.IsRunning = false
		chain.currIndex--
		chain.startRuntime = runtime + chain.tweens[chain.currIndex].duration
	}

	return tween.current
}

// region private

type tween struct {
	duration          float32
	from, to, current []float32
	easing            func(progress float32) float32
}

func (chain *Tween) last() *tween { return &chain.tweens[len(chain.tweens)-1] }

// endregion
