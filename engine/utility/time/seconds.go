package time

import (
	"fmt"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/flag"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/time/unit"
	"time"
)

//=================================================================
// setters

func SetScale(scale float32) { internal.TimeScale = scale }

//=================================================================
// getters

func AsClock24(seconds float32, divider string, units int) string {
	var ts = time.Duration(seconds * float32(time.Second))
	return formatTimeParts(ts, divider, units, false, false)
}
func AsClock12(seconds float32, divider string, units int, amPm bool) string {
	var ts = time.Duration(seconds * float32(time.Second))
	return formatTimeParts(ts, divider, units, true, amPm)
}

func Scale() float32                { return internal.TimeScale }
func Clock() float32                { return internal.Clock }
func FrameDelta() float32           { return internal.Delta }
func FrameRate() float32            { return internal.FrameRate }
func FrameRateAverage() float32     { return internal.FrameRateAverage }
func FrameCount() uint64            { return internal.FrameCount }
func Runtime() float32              { return internal.Runtime }
func RealFrameDelta() float32       { return internal.RealDelta }
func RealFrameRate() float32        { return internal.RealFrameRate }
func RealFrameRateAverage() float32 { return internal.RealFrameRateAverage }
func RealFrameCount() uint64        { return internal.RealFrameCount }
func RealRuntime() float32          { return internal.RealRuntime }

func ToMilliseconds(seconds float32) float32 { return seconds * 1000 }
func ToMinutes(secodns float32) float32      { return secodns / 60 }
func ToHours(seconds float32) float32        { return seconds / 3600 }
func ToDays(seconds float32) float32         { return seconds / 86400 }
func ToWeeks(seconds float32) float32        { return seconds / 604800 }

func FromMilliseconds(milliseconds float32) float32 { return milliseconds / 1000 }
func FromMinutes(minutes float32) float32           { return minutes * 60 }
func FromHours(hours float32) float32               { return hours * 3600 }
func FromDays(days float32) float32                 { return days * 86400 }
func FromWeeks(weeks float32) float32               { return weeks * 604800 }

//=================================================================
// private

func formatTimeParts(ts time.Duration, divider string, units int, is12Hour, amPm bool) string {
	var parts []string
	var counter = 0

	var conditionalSep = func() string {
		if counter > 0 {
			return divider
		}
		return ""
	}

	if flag.IsOn(units, unit.Day) {
		var val = int(ts.Hours() / 24)
		parts = append(parts, fmt.Sprintf("%02d", val))
		counter++
	}

	if flag.IsOn(units, unit.Hour) {
		var sep = conditionalSep()
		var val int
		if is12Hour {
			var h = int((ts % (24 * time.Hour)) / time.Hour)
			val = int(number.Wrap(float32(h), 12))
		} else {
			val = int((ts % (24 * time.Hour)) / time.Hour)
		}
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if flag.IsOn(units, unit.Minute) {
		var sep = conditionalSep()
		var val = int((ts % time.Hour) / time.Minute)
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if flag.IsOn(units, unit.Second) {
		var sep = conditionalSep()
		var val = int((ts % time.Minute) / time.Second)
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if flag.IsOn(units, unit.Millisecond) {
		var val = int((ts % time.Second) / time.Millisecond)
		var dot = ""
		if flag.IsOn(units, unit.Second) {
			dot = "."
		}
		var sep = ""
		if dot == "" && counter > 0 {
			sep = divider
		}
		parts = append(parts, sep+dot+fmt.Sprintf("%03d", val))
		counter++
	}

	if is12Hour && amPm {
		var sep = " "
		var amPm = "AM"
		if int(ts.Hours())%24 >= 12 {
			amPm = "PM"
		}
		parts = append(parts, sep+amPm)
	}

	return collection.ToText(parts, "")
}
