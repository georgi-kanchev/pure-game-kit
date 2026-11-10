package motion

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type Tween struct {
	TimeScale float32

	actions    []action
	index      int
	time       float32
	wasPlaying bool

	lastUpdateFrame, justFinishedFrame uint64
}

func NewTween(startingValues ...float32) *Tween {
	if len(startingValues) == 0 {
		return &Tween{}
	}

	var newChain = &Tween{actions: make([]action, 0), index: 1, TimeScale: 1}
	var itemsCopy = make([]float32, len(startingValues))
	copy(itemsCopy, startingValues)

	newChain.actions = append(newChain.actions, action{duration: 0, from: itemsCopy, to: itemsCopy, current: itemsCopy})
	return newChain
}

//=================================================================

func (t *Tween) ForceReset() {
	t.tryUpdate()
	t.actions = t.actions[:1] // leave only the starting values, clear everything else
	t.time, t.index = 0, 1
}
func (t *Tween) ForceFinish() {
	t.tryUpdate()
	var last = t.last()
	copy(t.actions[0].from, last.to)    // preserve last action values
	copy(t.actions[0].to, last.to)      // by setting them as starting values
	copy(t.actions[0].current, last.to) // as if from constructor
	t.ForceReset()
}

func (t *Tween) GoTo(duration float32, easing func(progress float32) float32, targets ...float32) *Tween {
	if len(t.actions) == 0 {
		return t
	}

	var lastTween, copyTo = t.last(), make([]float32, len(targets))
	var copyFrom, copyCurr = make([]float32, len(lastTween.to)), make([]float32, len(lastTween.to))
	copy(copyFrom, lastTween.to)
	copy(copyTo, targets)
	copy(copyCurr, lastTween.to)

	t.actions = append(t.actions, action{duration: duration, from: copyFrom, to: copyTo, current: copyCurr, easing: easing})
	t.tryUpdate()
	return t
}
func (t *Tween) Wait(seconds float32) *Tween {
	if len(t.actions) > 0 {
		t.GoTo(seconds, nil, t.last().to...)
	}
	t.tryUpdate()
	return t
}

func (t *Tween) CurrentValues() []float32 {
	t.tryUpdate()

	if !t.IsPlaying() {
		return t.last().current
	}

	var act = &t.actions[t.index]
	var actionDone = t.time > act.duration

	if len(t.actions) == 0 || len(act.from) != len(act.to) || len(act.from) != len(act.current) {
		return nil
	}

	for i := range act.current {
		if actionDone {
			act.current[i] = act.to[i]
			continue
		}
		var progress = t.time / act.duration
		var ease float32 = progress

		if act.easing != nil {
			ease = act.easing(progress)
		}

		act.current[i] = number.Map(ease, 0, 1, act.from[i], act.to[i])
	}

	if actionDone {
		t.time = 0
		t.index++
	}
	return act.current
}

func (t *Tween) IsPlaying() bool {
	t.tryUpdate()
	return t.index < len(t.actions)
}
func (t *Tween) IsJustFinished() bool {
	t.tryUpdate()
	return t.justFinishedFrame == internal.Frame
}

// =================================================================
// private

type action struct {
	duration          float32
	from, to, current []float32
	easing            func(progress float32) float32
}

func (t *Tween) last() *action { return &t.actions[len(t.actions)-1] }

func (t *Tween) tryUpdate() {
	if t.lastUpdateFrame == internal.Frame {
		return
	}
	t.lastUpdateFrame = internal.Frame

	if t.IsPlaying() {
		t.time += internal.FrameDelta * t.TimeScale
	} else if t.wasPlaying {
		t.ForceFinish()
		t.justFinishedFrame = internal.Frame
	}
	t.wasPlaying = t.IsPlaying()
}
