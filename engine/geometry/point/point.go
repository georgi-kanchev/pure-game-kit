package point

import (
	"math"
	"pure-kit/engine/utility/angle"
	"pure-kit/engine/utility/number"
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

	x -= float32(math.Mod(float64(x), float64(gridX)))
	y -= float32(math.Mod(float64(y), float64(gridY)))

	return x, y
}

func MoveAtAngle(x, y, angle, step float32) (float32, float32) {
	angle = number.Wrap(angle, 360)
	var rad = math.Pi / 180 * angle
	var dirX = float32(math.Cos(float64(rad)))
	var dirY = float32(math.Sin(float64(rad)))

	if dirX == 0 && dirY == 0 {
		return x, y
	}

	var length = float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
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
	var rad = float32(math.Pi/180) * angle
	var tx, ty = x - targetX, y - targetY
	var cosA = float32(math.Cos(float64(rad)))
	var sinA = float32(math.Sin(float64(rad)))
	var rx, ry = tx*cosA - ty*sinA, tx*sinA + ty*cosA

	return rx + targetX, ry + targetY
}

func DistanceToPoint(x, y, targetX, targetY float32) float32 {
	var dirX, dirY = targetX - x, targetY - y
	return float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
}
func DirectionToPoint(x, y, targetX, targetY float32) (dirX, dirY float32) {
	var length = DistanceToPoint(x, y, targetX, targetY)
	dirX, dirY = targetX-x, targetY-y
	return dirX / length, dirY / length
}
