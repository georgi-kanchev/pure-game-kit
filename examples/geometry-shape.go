package example

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/geometry/shape"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GeometryShape() {
	var cam = graphics.NewCamera(1)
	var shape = shape.New([2]float32{}, [2]float32{100, 100}, [2]float32{-100, 100})

	shape.ScaleY = 5
	shape.X += 200
	shape.Y -= 200

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		shape.Angle++
		var mx, my = cam.MousePosition()
		var col = condition.If(shape.Contains(mx, my), color.Red, color.Green)

		cam.DrawLinesPath(5, col, shape.Corners()...)
	}
}
