package internal

import (
	"pure-kit/engine/utility/number"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Clock, Delta, FrameRate, FrameRateAverage, Runtime float32
var RealDelta, RealFrameRate, RealFrameRateAverage, RealRuntime float32
var FrameCount, RealFrameCount uint64
var TimeScale float32 = 1

var CallAfter = make(map[float32][]func())
var CallFor = make(map[float32][]func(remaining float32))

func Update() {
	updateData()
	updateTimers()
	updateFlows()
	updateStates()
	updateKeysAndButtons()
}

//=================================================================
// private

const deltaMax float32 = 0.1

var prevClock float32

func updateData() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())

	Clock = secondsSinceMidnight

	if prevClock == 0 {
		prevClock = Clock
	}

	if prevClock > Clock { // we hit midnight
		prevClock = Clock - Delta
	}

	RealDelta = rl.GetFrameTime()
	RealRuntime += RealDelta
	RealFrameRate = 1.0 / RealDelta
	RealFrameRateAverage = float32(RealFrameCount) / RealRuntime
	RealFrameCount++

	Delta = number.Smallest(RealDelta*TimeScale, deltaMax)
	Runtime += Delta
	FrameRate = 1.0 / Delta
	FrameRateAverage = float32(FrameCount) / Runtime
	if RealDelta < deltaMax {
		FrameCount++
	}

	prevClock = Clock
}
