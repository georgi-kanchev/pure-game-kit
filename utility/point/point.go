package point

import (
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/direction"
	"pure-game-kit/utility/number"
)

func Snap(x, y, gridX, gridY float32) (float32, float32) {
	if gridX == 0 && gridY == 0 {
		return x, y
	}

	if x < 0 {
		y -= gridX
	}
	if y < 0 {
		y -= gridY
	}

	x -= number.DivisionRemainder(x, gridX)
	y -= number.DivisionRemainder(y, gridY)

	return x, y
}

func MoveAtAngle(x, y, angle, step float32) (float32, float32) {
	var dirX, dirY = dir(angle)
	if dirX == 0 && dirY == 0 {
		return x, y
	}

	var length = direction.Length(dirX, dirY)
	x += (dirX / length) * step
	y += (dirY / length) * step

	return x, y
}
func MoveToPoint(x, y, targetX, targetY float32, step float32) (float32, float32) {
	if x == targetX && y == targetY {
		return x, y
	}

	var angle = angle.BetweenPoints(x, y, targetX, targetY)
	x, y = MoveAtAngle(x, y, angle, step)

	if DistanceToPoint(x, y, targetX, targetY) <= step {
		return targetX, targetY
	}
	return x, y
}
func MoveByPercent(x, y, targetX, targetY float32, percent float32) (float32, float32) {
	x = number.Map(percent, 0, 100, x, targetX)
	y = number.Map(percent, 0, 100, y, targetY)
	return x, y
}

func RotateAroundPoint(x, y, targetX, targetY, angle float32) (float32, float32) {
	var rad = rad(angle)
	var tx, ty = x - targetX, y - targetY
	var cosA = number.Cosine(rad)
	var sinA = number.Sine(rad)
	var rx, ry = tx*cosA - ty*sinA, tx*sinA + ty*cosA

	return rx + targetX, ry + targetY
}

func DistanceToPoint(x, y, targetX, targetY float32) float32 {
	return direction.Length(targetX-x, targetY-y)
}
func DirectionToPoint(x, y, targetX, targetY float32) (dirX, dirY float32) {
	var length = DistanceToPoint(x, y, targetX, targetY)
	dirX, dirY = targetX-x, targetY-y
	return dirX / length, dirY / length
}

//=================================================================
// private

func dir(ang float32) (float32, float32) {
	return direction.FromAngle(ang)
}
func rad(ang float32) float32 {
	return angle.ToRadians(ang)
}
