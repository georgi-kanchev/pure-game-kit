package internal

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Clock, FrameDelta, Runtime float32
var Frame uint64
var FPS int32

var FixedStepCount int
var FixedTickDelta float32 = 1.0 / 60.0
var FixedStepAccumulatedTime float32

// private ========================================================

func UpdateTimeData() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())

	Clock = secondsSinceMidnight
	FrameDelta = min(rl.GetFrameTime(), 0.25)
	FPS = rl.GetFPS()
	Runtime += FrameDelta
	Frame++

	// custom fixed tick rate (see time package) ==================

	FixedStepAccumulatedTime += FrameDelta
	FixedStepCount = 0
	for FixedStepAccumulatedTime >= FixedTickDelta {
		FixedStepAccumulatedTime -= FixedTickDelta
		FixedStepCount++
	}
}
