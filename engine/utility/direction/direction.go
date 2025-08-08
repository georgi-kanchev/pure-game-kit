package direction

import (
	"math"
	"pure-kit/engine/utility/number"
)

func FromAngle(angle float32) (x, y float32) {
	var rad = number.Wrap(angle, 360) * (math.Pi / 180)
	x = float32(math.Cos(float64(rad)))
	y = float32(math.Sin(float64(rad)))
	return
}

func Dot(ax, ay, bx, by float32) float32 {
	return ax*bx + ay*by
}
func Normalize(x, y float32) (newX, newY float32) {
	var length = Length(x, y)
	if length == 0 {
		return 0, 0
	}
	newX = x / length
	newY = y / length
	return
}
func Length(x, y float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y)))
}
