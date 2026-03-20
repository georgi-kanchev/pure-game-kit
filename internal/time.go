package internal

import (
	"pure-game-kit/utility/number"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Clock, DeltaTime, FPS, AverageFPS, FrameTime, Runtime float32
var FrameCount uint64

var CallAfter = make(map[float32][]func())
var CallFor = make(map[float32][]func(remaining float32))

var FrameStart time.Time

func Update() {
	if FrameCount == 0 {
		initData()
	}

	updateTimeData()
	updateTimers()
	updateInput()
	updateMusic()
	updateScreens()
}

//=================================================================
// private

const deltaMax, alpha float32 = 0.1, 0.1

var prevClock, smoothDelta float32

func updateTimeData() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = float32(now.Sub(midnight).Seconds())

	Clock = secondsSinceMidnight

	if prevClock == 0 {
		prevClock = Clock
	}

	if prevClock > Clock { // we hit midnight
		prevClock = Clock - DeltaTime
	}

	DeltaTime = number.Smallest(rl.GetFrameTime(), deltaMax)
	Runtime += DeltaTime
	FrameCount++
	FPS = 1.0 / DeltaTime

	if smoothDelta == 0 {
		smoothDelta = DeltaTime
	}

	smoothDelta = (DeltaTime * alpha) + (smoothDelta * (1.0 - alpha))
	AverageFPS = 1.0 / smoothDelta

	prevClock = Clock
}
