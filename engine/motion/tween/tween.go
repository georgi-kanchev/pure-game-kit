package tween

// type Chain struct {
// 	tweens    []tween
// 	currIndex int
// 	elapsed   float32
// 	playing   bool
// }

// type tween struct {
// 	duration             float32
// 	from, to, current    []float32
// 	easing               func(progress float32) float32
// 	repeats, repeatsLeft int
// 	whenThere            func()
// 	whileGoing           func(progress float32, current []float32)
// }

// func From(items []float32) *Chain {
// 	if len(items) == 0 {
// 		return &Chain{}
// 	}

// 	var newChain = &Chain{
// 		tweens:  make([]tween, 0),
// 		elapsed: 0.0,
// 	}

// 	var itemsCopy = make([]float32, len(items))
// 	copy(itemsCopy, items)

// 	var newValue = tween{
// 		duration:    0,
// 		from:        itemsCopy,
// 		to:          itemsCopy,
// 		current:     itemsCopy,
// 		repeats:     0,
// 		repeatsLeft: 0,
// 	}
// 	newChain.tweens = append(newChain.tweens, newValue)

// 	return newChain
// }
// func (chain *Chain) GoTo(targets []float32, duration float32, easing func(progress float32) float32) *Chain {
// 	if len(chain.tweens) == 0 {
// 		return chain
// 	}

// 	var lastTween = chain.last()
// 	var copyFrom = make([]float32, len(lastTween.to))
// 	var copyTo = make([]float32, len(targets))
// 	var copyCurr = make([]float32, len(lastTween.to))
// 	copy(copyFrom, lastTween.to)
// 	copy(copyTo, targets)
// 	copy(copyCurr, lastTween.to)

// 	var newTween = tween{
// 		duration:    duration,
// 		from:        copyFrom,
// 		to:          copyTo,
// 		current:     copyCurr,
// 		easing:      easing,
// 		repeats:     1,
// 		repeatsLeft: 1,
// 	}
// 	chain.tweens = append(chain.tweens, newTween)

// 	return chain
// }
// func (chain *Chain) GoBack() *Chain {
// 	if len(chain.tweens) > 0 {
// 		var lastTween = chain.last()
// 		chain.GoTo(lastTween.from, lastTween.duration, lastTween.easing)
// 	}
// 	return chain
// }
// func (chain *Chain) Wait(seconds float32) *Chain {
// 	if len(chain.tweens) > 0 {
// 		var lastTween = chain.last()
// 		chain.GoTo(lastTween.to, seconds, nil)
// 	}
// 	return chain
// }
// func (chain *Chain) Repeat(times int) *Chain {
// 	for i := len(chain.tweens) - 1; i >= 0; i-- {
// 		chain.tweens[i].repeatsLeft += times
// 		chain.tweens[i].repeats += times
// 	}
// 	return chain
// }

// func (chain *Chain) CallWhenDone(function func()) *Chain {
// 	if len(chain.tweens) > 0 {
// 		chain.last().whenThere = function
// 	}
// 	return chain
// }
// func (chain *Chain) CallWhileDoing(function func(progress float32, current []float32)) *Chain {
// 	if len(chain.tweens) > 0 {
// 		chain.last().whileGoing = function
// 	}
// 	return chain
// }

// func (chain *Chain) Restart() {
// 	chain.currIndex = 0
// 	chain.elapsed = 0
// 	chain.playing = true

// 	for i := range chain.tweens {
// 		var tween = &chain.tweens[i]
// 		if i == 0 {
// 			copy(tween.current, tween.from)
// 			continue
// 		}

// 		tween.repeatsLeft = tween.repeats
// 		copy(tween.current, chain.tweens[i-1].to)

// 	}
// }
// func (chain *Chain) Pause(paused bool) { chain.playing = !paused }
// func (chain *Chain) Advance(deltaTime float32) []float32 {
// 	if len(chain.tweens) == 0 {
// 		return []float32{}
// 	}

// 	var tween = chain.current()

// 	if !chain.playing {
// 		return tween.current
// 	}

// 	if len(tween.from) != len(tween.to) || len(tween.from) != len(tween.current) {
// 		return []float32{}
// 	}

// 	chain.elapsed += deltaTime

// 	var tweenDone = chain.elapsed > tween.duration

// 	for i := range tween.current {
// 		if tweenDone {
// 			tween.current[i] = tween.to[i]
// 			continue
// 		}
// 		var progress = chain.elapsed / tween.duration
// 		var ease float32 = progress

// 		if tween.easing != nil {
// 			ease = tween.easing(progress)
// 		}

// 		tween.current[i] = number.Map(ease, 0, 1, tween.from[i], tween.to[i])
// 	}

// 	if tweenDone {
// 		tween.repeatsLeft--
// 		chain.elapsed = 0
// 		chain.currIndex++

// 		if tween.whenThere != nil {
// 			tween.whenThere()
// 		}
// 	} else if tween.whileGoing != nil {
// 		tween.whileGoing(chain.elapsed/tween.duration, tween.current)
// 	}

// 	var chainDone = chain.currIndex >= int(len(chain.tweens))

// 	if chainDone {
// 		if chain.tweens[0].repeatsLeft > 0 {
// 			chain.currIndex = 0
// 			chain.elapsed = 0
// 			return tween.current
// 		}

// 		chain.playing = false
// 		chain.currIndex--
// 		chain.elapsed = chain.tweens[chain.currIndex].duration
// 	} else if chain.current().repeatsLeft < tween.repeatsLeft {
// 		chain.currIndex = 0
// 		chain.elapsed = 0
// 		return tween.current
// 	}

// 	return tween.current
// }

// // region private

// func (chain *Chain) last() *tween    { return &chain.tweens[len(chain.tweens)-1] }
// func (chain *Chain) current() *tween { return &chain.tweens[chain.currIndex] }

// // endregion
