package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Shapes() {
	var cam = graphics.NewCamera(1)
	var shape = geometry.NewShapeCorners([][2]float32{
		{0, 0},
		{50, -20},
		{100, 0},
		{0, 100},
		{50, 120},
		{100, 100},
	}...)
	var triangle = geometry.NewShapeCorners([2]float32{}, [2]float32{100, 100}, [2]float32{-100, 100})
	var rectangle = geometry.NewShapeRectangle(700, 500, 0.5, 0.5)
	var circle = geometry.NewShapeSides(500, 16)

	shape.ScaleX, shape.ScaleY = 5, 5
	shape.X += 180
	shape.Y -= 200
	rectangle.Angle = 45

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		shape.Angle += time.FrameDelta() * 60
		var mx, my = cam.MousePosition()
		var colShape = condition.If(shape.IsOverlappingShapes(triangle), color.Red, color.Green)
		var colRect = condition.If(rectangle.IsCrossingShapes(shape), color.Brown, color.Cyan)
		var colCircle = condition.If(circle.IsContainingShapes(triangle), color.Yellow, color.Pink)

		triangle.X, triangle.Y = mx, my

		var crossPoints = circle.CrossPointsWithShapes(shape)

		cam.DrawLinesPath(8, colRect, rectangle.CornerPoints()...)
		cam.DrawLinesPath(8, colCircle, circle.CornerPoints()...)
		cam.DrawLinesPath(8, colShape, shape.CornerPoints()...)
		cam.DrawLinesPath(8, color.White, triangle.CornerPoints()...)

		for _, v := range crossPoints {
			cam.DrawCircle(v[0], v[1], 16, color.Green)
		}
	}
}
