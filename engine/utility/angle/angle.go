package utility

import "math"

type Angle float32

const (
	DegreesRight     Angle = 0
	DegreesDownRight Angle = 45
	DegreesDown      Angle = 90
	DegreesDownLeft  Angle = 135
	DegreesLeft      Angle = 180
	DegreesUpLeft    Angle = 225
	DegreesUp        Angle = 270
	DegreesUpRight   Angle = 315
)

func DegreesToRadians(degrees float32) float32 {
	degrees = wrap(degrees, 360)
	return degrees * (math.Pi / 180)
}
func RadiansToDegrees(radians float32) float32 {
	deg := radians * (180 / math.Pi)
	return wrap(deg, 360)
}

func Rotate(angle, target, speed float32) float32 {
	if speed < 0 {
		speed = -speed
	}

	// wrap both angles to [0, 360)
	angle = wrap(angle, 360)
	target = wrap(target, 360)

	// Shortest signed difference from angle to target
	diff := wrap(target-angle+180, 360) - 180
	checkedSpeed := speed

	// Snap if within rotation threshold or just over 360 wraparound
	if math.Abs(float64(diff)) < float64(checkedSpeed) || math.Abs(float64(diff)) > float64(360-checkedSpeed) {
		return target
	}

	// Rotate in shortest direction
	if diff > 0 {
		return wrap(angle+checkedSpeed, 360)
	}
	return wrap(angle-checkedSpeed, 360)
}
func Face(angle, target, progress float32) float32 {
	angle = wrap(angle, 360)
	target = wrap(target, 360)

	diff := target - angle
	if diff <= -180 {
		diff += 360
	} else if diff > 180 {
		diff -= 360
	}

	interpolated := angle + diff*progress
	return wrap(interpolated, 360)
}

func IsBehind(angle, target float32) bool {
	diff := wrap(angle-target, 360)
	switch {
	case diff >= 0 && diff < 180:
		return false
	case diff >= -180 && diff < 0:
		return true
	case diff >= -360 && diff < -180:
		return false
	case diff >= 180 && diff < 360:
		return true
	default:
		return false
	}
}
func IsWithin(angle, lower, upper float32) bool {
	if !IsBehind(angle, lower) || IsBehind(angle, upper) {
		return false
	}
	return Distance(angle, lower) < 180 && Distance(angle, upper) < 180
}

func Distance(angle, target float32) float32 {
	diff := wrap(target-angle, 360)
	if diff < -180 {
		diff += 360
	} else if diff >= 180 {
		diff -= 360
	}
	return float32(math.Abs(float64(diff)))
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
func Dot(angle, target float32) float32 {
	ax := float32(math.Cos(float64(angle) * math.Pi / 180))
	ay := float32(math.Sin(float64(angle) * math.Pi / 180))
	bx := float32(math.Cos(float64(target) * math.Pi / 180))
	by := float32(math.Sin(float64(target) * math.Pi / 180))
	return ax*bx + ay*by
}
func Reflect(angle, surfaceAngle float32) float32 {
	return wrap(2*surfaceAngle-angle+180, 360)
}
func Reverse(angle float32) float32 {
	return wrap(angle-180, 360)
}

func BetweenPoints(x, y, targetX, targetY float32) float32 {
	dx := targetX - x
	dy := targetY - y
	angleRad := math.Atan2(float64(dy), float64(dx))
	angleDeg := float32(angleRad * 180 / math.Pi)
	return wrap(angleDeg, 360)
}

func ToDirection(angle float32) (x, y float32) {
	radians := float64(angle) * math.Pi / 180
	x = float32(math.Cos(radians))
	y = float32(math.Sin(radians))
	return
}

// region private
func wrap(number, target float32) float32 { // copied from utility/number
	if target == 0 {
		return 0
	}
	return float32(math.Mod(math.Mod(float64(number), float64(target))+float64(target), float64(target)))
}

// endregion
