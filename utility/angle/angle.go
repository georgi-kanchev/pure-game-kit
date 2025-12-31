// Helper functions for an angle in 360 degrees. Useful for working with 2D objects.
package angle

import (
	"math"
	"pure-game-kit/utility/number"
)

func ToRadians(degrees float32) float32 {
	return degrees * (math.Pi / 180) // not wrapped for performance
}
func FromRadians(radians float32) float32 {
	return radians * (180 / math.Pi) // not wrapped for performance
}

func Rotate(angle, target, speed float32) float32 {
	if speed < 0 {
		speed = -speed
	}

	var diff = number.Wrap(target-angle+180, 0, 360) - 180 // only wrap the delta, not both angles

	if number.Unsign(diff) <= speed {
		return target
	}

	if diff > 0 {
		return angle + speed
	}
	return angle - speed
}
func Face(angle, target, progress float32) float32 {
	var diff = target - angle

	if diff <= -180 {
		diff += 360
	} else if diff > 180 {
		diff -= 360
	}

	return angle + diff*progress
}

func IsBehind(angle, target float32) bool {
	var diff = number.Wrap(angle-target, -180, 180)
	return diff < 0
}
func IsWithin(angle, lower, upper float32) bool {
	if !IsBehind(angle, lower) || IsBehind(angle, upper) {
		return false
	}
	return Distance(angle, lower) < 180 && Distance(angle, upper) < 180
}

func Distance(angle, target float32) float32 {
	var diff = number.Wrap(target-angle, -180, 180)
	if diff < 0 {
		return -diff
	}
	return diff
}
func Limit(angle, lower, upper float32) float32 {
	if !IsWithin(angle, lower, upper) {
		if Distance(angle, lower) < Distance(angle, upper) {
			return lower
		}
		return upper
	}
	return angle
}
func Reflect(angle, surfaceAngle float32) float32 {
	return 2*surfaceAngle - angle
}
func Reverse(angle float32) float32 {
	return angle - 180
}

func BetweenPoints(x, y, targetX, targetY float32) float32 {
	return FromRadians(float32(math.Atan2(float64(targetY-y), float64(targetX-x))))
}

func FromDirection(x, y float32) float32 {
	return FromRadians(float32(math.Atan2(float64(y), float64(x))))
}
