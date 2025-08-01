package internal

import (
	"math"
	"time"
)

const deltaMax float32 = 0.1

var prevClock float32

var Clock, Delta, FrameRate, FrameRateAverage, Runtime float32
var RealDelta, RealFrameRate, RealFrameRateAverage, RealRuntime float32
var FrameCount, RealFrameCount uint64
var IsPaused bool

var CallAfter = make(map[float32][]func())
var CallFor = make(map[float32][]func(remaining float32))

func Update() {
	updateData()
	updateTimers()
	updateFlows()
	updateStates()
	updateKeysAndButtons()
}

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

	RealDelta = Clock - prevClock
	RealRuntime += RealDelta
	RealFrameRate = 1.0 / RealDelta
	RealFrameRateAverage = float32(RealFrameCount) / RealRuntime
	RealFrameCount++

	Delta = 0
	FrameRate = 0
	if !IsPaused {
		Delta = float32(math.Min(float64(RealDelta), float64(deltaMax)))
		Runtime += Delta
		FrameRate = 1.0 / Delta
		FrameRateAverage = float32(FrameCount) / Runtime
		if RealDelta < deltaMax {
			FrameCount++
		}
	}

	prevClock = Clock
}
