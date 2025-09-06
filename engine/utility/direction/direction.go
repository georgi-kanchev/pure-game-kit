package direction

import (
	"pure-kit/engine/utility/angle"
	"pure-kit/engine/utility/number"
)

func FromAngle(angle float32) (x, y float32) {
	var rad = rad(angle)
	return number.Cosine(rad), number.Sine(rad)
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
	return number.SquareRoot(x*x + y*y)
}

//=================================================================
// private

func rad(ang float32) float32 {
	return angle.ToRadians(ang)
}
