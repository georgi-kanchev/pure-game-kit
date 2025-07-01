package internal

import (
	"math"
	"time"
)

const deltaMax float64 = 0.1

var prevClock float64 = 0.0

var Clock, Delta, FrameRate, FrameRateAverage, Runtime float64
var RealDelta, RealFrameRate, RealFrameRateAverage, RealRuntime float64
var FrameCount, RealFrameCount uint64

func Update() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = now.Sub(midnight).Seconds()

	Clock = secondsSinceMidnight
	if prevClock == 0 {
		prevClock = Clock
	}

	RealDelta = Clock - prevClock
	RealRuntime += RealDelta
	RealFrameRate = 1.0 / RealDelta
	RealFrameRateAverage = float64(RealFrameCount) / RealRuntime
	RealFrameCount++

	Delta = math.Min(RealDelta, deltaMax)
	Runtime += Delta
	FrameRate = 1.0 / Delta
	FrameRateAverage = float64(FrameCount) / Runtime
	if RealDelta < Delta {
		FrameCount++
	}

	prevClock = Clock
}
