package internal

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Clock, FrameDelta, Runtime float32
var Frame uint64
var FPS int32

// private ========================================================

const deltaMax, alpha float32 = 0.1, 0.1

var prev time.Time = time.Now()

func UpdateTimeData() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())

	Clock = secondsSinceMidnight
	FrameDelta = rl.GetFrameTime()
	FPS = rl.GetFPS()
	Runtime += FrameDelta
	Frame++

	prev = now
}
