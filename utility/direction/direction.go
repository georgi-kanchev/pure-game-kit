/*
Helper functions for a directional vector. Similar angle but interpreting data differently.
Useful for working with 2D objects.
*/
package direction

import (
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/number"
)

func FromAngle(angle float32) (x, y float32) {
	var rad = rad(angle)
	return number.Cosine(rad), number.Sine(rad)
}

func Dot(ax, ay, bx, by float32) float32 {
	return ax*bx + ay*by
}
func Normalize(x, y float32) (float32, float32) {
	var length = Length(x, y)
	if length == 0 {
		return 0, 0
	}
	return x / length, y / length
}
func Length(x, y float32) float32 {
	return number.SquareRoot(x*x + y*y)
}
func BetweenPoints(x, y, targetX, targetY float32) (float32, float32) {
	return Normalize(targetX-x, targetY-y)
}

//=================================================================
// private

func rad(ang float32) float32 {
	return angle.ToRadians(ang)
}
