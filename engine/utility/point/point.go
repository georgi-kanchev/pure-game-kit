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
		x -= gridX
	}
	if y < 0 {
		y -= gridY
	}

	x -= float32(math.Mod(float64(x), float64(gridX)))
	y -= float32(math.Mod(float64(y), float64(gridY)))

	return x, y
}

func MoveIn(x, y, dirX, dirY, distance float32) (float32, float32) {
	if dirX == 0 && dirY == 0 {
		return x, y
	}

	var m = float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
	var normX = dirX / m
	var normY = dirY / m
	var newX = x + normX*distance
	var newY = y + normY*distance

	return newX, newY
}
func MoveAt(x, y, angle, distance float32) (float32, float32) {
	angle = number.Wrap(angle, 360)
	var rad = math.Pi / 180 * angle
	var dirX = float32(math.Cos(float64(rad)))
	var dirY = float32(math.Sin(float64(rad)))

	return MoveIn(x, y, dirX, dirY, distance)
}
func MoveTo(x, y, targetX, targetY, distance float32) (float32, float32) {
	if x == targetX && y == targetY {
		return x, y
	}

	var dirX = targetX - x
	var dirY = targetY - y
	var resultX, resultY = MoveIn(x, y, dirX, dirY, distance)

	if Distance(resultX, resultY, targetX, targetY) < distance*0.51 {
		return targetX, targetY
	}
	return resultX, resultY
}
func MoveBy(x, y, targetX, targetY, percent float32) (float32, float32) {
	var newX = number.Map(percent, 0, 100, x, targetX)
	var newY = number.Map(percent, 0, 100, y, targetY)
	return newX, newY
}

func Distance(x1, y1, x2, y2 float32) float32 {
	var dx = x2 - x1
	var dy = y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func Direction(x1, y1, x2, y2 float32) (float32, float32) {
	var dirX = x2 - x1
	var dirY = y2 - y1
	var m = float32(math.Sqrt(float64(dirX*dirX + dirY*dirY)))
	return dirX / m, dirY / m
}
