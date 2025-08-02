package point

import (
	"math"
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

func MoveIn(x, y, directionX, directionY, distance float32) (float32, float32) {
	if directionX == 0 && directionY == 0 {
		return x, y
	}

	var length = float32(math.Sqrt(float64(directionX*directionX + directionY*directionY)))
	x += (directionX / length) * distance
	y += (directionY / length) * distance
	return x, y
}
func MoveAt(x, y, angle, distance float32) (float32, float32) {
	angle = number.Wrap(angle, 360)
	var rad = math.Pi / 180 * angle
	var dirX = float32(math.Cos(float64(rad)))
	var dirY = float32(math.Sin(float64(rad)))
	return MoveIn(x, y, dirX, dirY, distance)
}
func MoveTo(x, y, targetX, targetY float32, distance float32) (float32, float32) {
	if x == targetX && y == targetY {
		return x, y
	}

	x, y = MoveIn(x, y, targetX-x, targetY-y, distance)

	if Distance(x, y, targetX, targetY) <= distance {
		return targetX, targetY
	}
	return x, y
}
func MoveBy(x, y, targetX, targetY float32, percent float32) (float32, float32) {
	x = number.Map(percent, 0, 100, x, targetX)
	y = number.Map(percent, 0, 100, y, targetY)
	return x, y
}

func Distance(x, y, targetX, targetY float32) float32 {
	var dirX, dirY = targetX - x, targetY - y
	return float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
}
func Direction(x, y, targetX, targetY float32) (dirX, dirY float32) {
	var length = Distance(x, y, targetX, targetY)
	dirX, dirY = targetX-x, targetY-y
	return dirX / length, dirY / length
}
