package internal

import (
	"pure-game-kit/utility/number"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var Clock, DeltaTime, FrameRate, FrameRateAverage, Runtime float32
var RealDeltaTime, RealFrameRate, RealFrameRateAverage, RealRuntime float32
var FrameCount, RealFrameCount uint64
var TimeScale float32 = 1

var CallAfter = make(map[float32][]func())
var CallFor = make(map[float32][]func(remaining float32))

func Update() {
	updateTimeData()
	updateTimers()
	updateInput()
	updateMusic()
	updateAnimatedTiles()
	updateScreens()
}

//=================================================================
// private

const deltaMax float32 = 0.1

var prevClock float32

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

	RealDeltaTime = rl.GetFrameTime()
	RealRuntime += RealDeltaTime
	RealFrameRate = 1.0 / RealDeltaTime
	RealFrameRateAverage = float32(RealFrameCount) / RealRuntime
	RealFrameCount++

	DeltaTime = number.Smallest(RealDeltaTime*TimeScale, deltaMax)
	Runtime += DeltaTime
	FrameRate = 1.0 / DeltaTime
	FrameRateAverage = float32(FrameCount) / Runtime
	if RealDeltaTime < deltaMax {
		FrameCount++
	}

	prevClock = Clock
}
