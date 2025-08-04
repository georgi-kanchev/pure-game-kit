package shape

import (
	"math"
	"pure-kit/engine/utility/angle"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape struct {
	X, Y, Angle    float32
	ScaleX, ScaleY float32
	corners        [][2]float32
}

func New(corners ...[2]float32) Shape {
	if len(corners) == 0 {
		return Shape{}
	}
	return Shape{ScaleX: 1, ScaleY: 1, corners: append(corners, corners[0])}
}

func (shape *Shape) Corners() [][2]float32 {
	var result = make([][2]float32, len(shape.corners))
	for i := range shape.corners {
		var x, y = shape.corners[i][0], shape.corners[i][1]

		x *= shape.ScaleX
		y *= shape.ScaleY

		var rad = angle.ToRadians(shape.Angle)
		var sin, cos = float32(math.Sin(float64(rad))), float32(math.Cos(float64(rad)))
		var resultX = shape.X + (x*cos - y*sin)
		var resultY = shape.Y + (x*sin + y*cos)

		result[i] = [2]float32{resultX, resultY}
	}
	return result
}

func (shape *Shape) Contains(x, y float32) bool {
	var result = make([]rl.Vector2, len(shape.corners))
	var corners = shape.Corners()
	for i, p := range corners {
		result[i] = rl.Vector2{X: p[0], Y: p[1]}
	}

	return rl.CheckCollisionPointPoly(rl.Vector2{X: x, Y: y}, result)
}
