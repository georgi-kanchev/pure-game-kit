// Helper functions for time. Provides time conversions. As well as runtime and real time stats.
// Provides a way to slow down or speed up those runtime stats. All values are in seconds, unless
// a unit is explicitly specified otherwise. Useful for when dealing with motion, debugging or displaying
// timers/clocks etc.
package time

import (
	"pure-game-kit/packages/internal"
)

func Running() float32 { return internal.Runtime }
func Clock() float32   { return internal.Clock }
func Delta() float32   { return internal.FrameDelta }
func Frame() uint64    { return internal.Frame }
func FPS() int         { return int(internal.FPS) }

func ToMilliseconds(seconds float32) float32 { return seconds * 1000 }
func ToMinutes(secodns float32) float32      { return secodns / 60 }
func ToHours(seconds float32) float32        { return seconds / 3600 }

func FromMilliseconds(milliseconds float32) float32 { return milliseconds / 1000 }
func FromMinutes(minutes float32) float32           { return minutes * 60 }
func FromHours(hours float32) float32               { return hours * 3600 }

func AsTimer(seconds float32) (hr, min, sec, ms int) {
	var totalMs = int64(seconds * 1_000)
	hr = int(totalMs / 3_600_000)
	min = int((totalMs % 3_600_000) / 60_000)
	sec = int((totalMs % 60_000) / 1_000)
	ms = int(totalMs % 1_000)
	return hr, min, sec, ms
}
func AsClock24(seconds float32) (hr, min, sec int) {
	const dayInSeconds = 86400
	var totalSec = int(seconds) % dayInSeconds // wrap seconds in a 24-hour period (handles negative values as well)
	if totalSec < 0 {
		totalSec += dayInSeconds
	}

	hr = totalSec / 3600
	min = (totalSec % 3600) / 60
	sec = totalSec % 60
	return hr, min, sec
}
func AsClock12(seconds float32) (hr, min, sec int, pm bool) {
	var hr24 int
	hr24, min, sec = AsClock24(seconds)

	pm = hr24 >= 12
	hr = hr24 % 12
	if hr == 0 {
		hr = 12 // 00:00 is 12 AM, 12:00 is 12 PM
	}
	return hr, min, sec, pm
}
