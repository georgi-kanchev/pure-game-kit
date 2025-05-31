package motion

import (
	"math"
)

type Chain struct {
	tweens    []tween
	currIndex int32
	elapsed   float32
	playing   bool
}

type tween struct {
	duration             float32
	from, to, current    []float32
	easing               func(float32) float32
	repeats, repeatsLeft int32
}

func From(items []float32) *Chain {
	if len(items) == 0 {
		return &Chain{}
	}

	var newChain = &Chain{
		tweens:  make([]tween, 0),
		elapsed: 0.0,
	}

	var itemsCopy = make([]float32, len(items))
	copy(itemsCopy, items)

	var newValue = tween{
		duration:    0,
		from:        itemsCopy,
		to:          itemsCopy,
		current:     itemsCopy,
		repeats:     0,
		repeatsLeft: 0,
	}
	newChain.tweens = append(newChain.tweens, newValue)

	return newChain
}
func (chain *Chain) To(targets []float32, duration float32, easing func(float32) float32) *Chain {
	if len(chain.tweens) == 0 {
		return chain
	}

	var valuePrev = chain.tweens[len(chain.tweens)-1]
	var copyFrom = make([]float32, len(valuePrev.to))
	var copyTo = make([]float32, len(targets))
	var copyCurr = make([]float32, len(valuePrev.to))
	copy(copyFrom, valuePrev.to)
	copy(copyTo, targets)
	copy(copyCurr, valuePrev.to)

	var newTween = tween{
		duration:    duration,
		from:        copyFrom,
		to:          copyTo,
		current:     copyCurr,
		easing:      easing,
		repeats:     1,
		repeatsLeft: 1,
	}
	chain.tweens = append(chain.tweens, newTween)

	return chain
}
func (chain *Chain) Repeat(times int32) *Chain {
	for i := len(chain.tweens) - 1; i >= 0; i-- {
		chain.tweens[i].repeatsLeft += times
		chain.tweens[i].repeats += times
	}
	return chain
}

func (chain *Chain) Restart() {
	chain.currIndex = 0
	chain.elapsed = 0
	chain.playing = true

	for i := range chain.tweens {
		if i == 0 {
			copy(chain.tweens[i].current, chain.tweens[i].from)
			continue
		}

		chain.tweens[i].repeatsLeft = chain.tweens[i].repeats
		copy(chain.tweens[i].current, chain.tweens[i-1].to)

	}
}
func (chain *Chain) Pause(paused bool) {
	chain.playing = !paused
}

func (chain *Chain) Update(deltaTime float32) []float32 {
	if len(chain.tweens) == 0 {
		return []float32{}
	}

	var tween = &chain.tweens[chain.currIndex]

	if !chain.playing {
		return tween.current
	}

	if len(tween.from) != len(tween.to) || len(tween.from) != len(tween.current) {
		return []float32{}
	}

	chain.elapsed += deltaTime

	for i := range tween.current {
		if chain.elapsed > tween.duration {
			tween.current[i] = tween.to[i]
			continue
		}
		var ease = tween.easing(chain.elapsed / tween.duration)
		tween.current[i] = mapFloat(ease, 0, 1, tween.from[i], tween.to[i])
	}
	if chain.elapsed > tween.duration {
		tween.repeatsLeft--
		chain.elapsed = 0
		chain.currIndex++
	}

	if chain.currIndex >= int32(len(chain.tweens)) {
		if chain.tweens[0].repeatsLeft > 0 {
			chain.currIndex = 0
			chain.elapsed = 0
			return tween.current
		}

		chain.playing = false
		chain.currIndex--
		chain.elapsed = chain.tweens[chain.currIndex].duration
	} else if chain.tweens[chain.currIndex].repeatsLeft < tween.repeatsLeft {
		chain.currIndex = 0
		chain.elapsed = 0
		return tween.current
	}

	return tween.current
}

// region private

func mapFloat(number float32, fromA, fromB, toA, toB float32) float32 { // copied from utility/number
	if math.Abs(float64(fromB-fromA)) < 0.001 {
		return (toA + toB) / 2
	}
	value := ((number-fromA)/(fromB-fromA))*(toB-toA) + toA
	if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
		return toA
	}
	return value
}

// endregion
