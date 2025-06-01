package time

import (
	"fmt"
	"math"
	"pure-kit/engine/utility/number"
	"strings"
	"time"
)

type Unit int

const (
	Day Unit = 1 << iota
	Hour
	Minute
	Second
	Millisecond
)

type Conversion int

const deltaMax float64 = 0.1

var Clock, prevClock, Delta, DeltaRaw, FrameRate, FrameRateAverage float64
var FrameCount uint64
var Runtime float64

func Update() {
	var now = time.Now()
	var midnight = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var secondsSinceMidnight = now.Sub(midnight).Seconds()

	Clock = secondsSinceMidnight
	DeltaRaw = Clock - prevClock
	Delta = math.Min(DeltaRaw, deltaMax)
	Runtime += Delta
	FrameRate = 1.0 / Delta
	FrameRateAverage = float64(FrameCount) / Runtime
	FrameCount++
	prevClock = Clock
}

func AsClock24(seconds float64, divider string, units Unit) string {
	ts := time.Duration(seconds * float64(time.Second))
	return formatTimeParts(ts, divider, units, false, false)
}
func AsClock12(seconds float64, divider string, units Unit, AM_PM bool) string {
	ts := time.Duration(seconds * float64(time.Second))
	return formatTimeParts(ts, divider, units, true, AM_PM)
}

func SecondsToMilliseconds(seconds float64) float64 { return seconds * 1000 }
func SecondsToMinutes(secodns float64) float64      { return secodns / 60 }
func SecondsToHours(seconds float64) float64        { return seconds / 3600 }
func SecondsToDays(seconds float64) float64         { return seconds / 86400 }
func SecondsToWeeks(seconds float64) float64        { return seconds / 604800 }

func SecondsFromMilliseconds(milliseconds float64) float64 { return milliseconds / 1000 }
func SecondsFromMinutes(minutes float64) float64           { return minutes * 60 }
func SecondsFromHours(hours float64) float64               { return hours * 3600 }
func SecondsFromDays(days float64) float64                 { return days * 86400 }
func SecondsFromWeeks(weeks float64) float64               { return weeks * 604800 }

// region private

func formatTimeParts(ts time.Duration, divider string, units Unit, is12Hour, amPm bool) string {
	var parts []string
	counter := 0

	conditionalSep := func() string {
		if counter > 0 {
			return divider
		}
		return ""
	}

	if units&Day != 0 {
		val := int(ts.Hours() / 24)
		parts = append(parts, fmt.Sprintf("%02d", val))
		counter++
	}

	if units&Hour != 0 {
		sep := conditionalSep()
		var val int
		if is12Hour {
			h := int((ts % (24 * time.Hour)) / time.Hour)
			val = int(number.Wrap(float32(h), 12))
		} else {
			val = int((ts % (24 * time.Hour)) / time.Hour)
		}
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if units&Minute != 0 {
		sep := conditionalSep()
		val := int((ts % time.Hour) / time.Minute)
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if units&Second != 0 {
		sep := conditionalSep()
		val := int((ts % time.Minute) / time.Second)
		parts = append(parts, sep+fmt.Sprintf("%02d", val))
		counter++
	}

	if units&Millisecond != 0 {
		val := int((ts % time.Second) / time.Millisecond)
		dot := ""
		if units&Second != 0 {
			dot = "."
		}
		sep := ""
		if dot == "" && counter > 0 {
			sep = divider
		}
		parts = append(parts, sep+dot+fmt.Sprintf("%03d", val))
		counter++
	}

	if is12Hour && amPm {
		sep := " "
		amPm := "AM"
		if int(ts.Hours())%24 >= 12 {
			amPm = "PM"
		}
		parts = append(parts, sep+amPm)
	}

	return strings.Join(parts, "")
}

// endregion
