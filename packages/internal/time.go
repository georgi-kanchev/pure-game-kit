package internal

import (
	"time"
)

var Clock, TickDelta, FrameDelta, TPS, FPS, Runtime float32
var TickBusy float32
var TargetTPS uint16

// private ========================================================

const deltaMax, alpha float32 = 0.1, 0.1

var prev time.Time = time.Now()

func updateTimeData() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())

	Clock = secondsSinceMidnight

	TickDelta = float32(time.Since(prev).Seconds())
	TPS = 1.0 / TickDelta
	Runtime += TickDelta

	prev = now
}
