package example

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/window"
)

func ShapesGrids() {
	var cam = graphics.NewCamera(4)
	var grid = geometry.NewShapeGrid(32, 32)
	var shape = geometry.NewShapeCorners(
		[2]float32{},
		[2]float32{50, -20},
		[2]float32{100, 0},
		[2]float32{0, 100},
		[2]float32{50, 120},
		[2]float32{100, 100})

	for i := -8; i < 8; i++ {
		for j := -8; j < 8; j++ {
			grid.SetShapesAtCell(i, j, geometry.NewShapeRectangle(24, 24, 0.5, 0.5))
		}
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 32, 32, color.Red)

		var mx, my = cam.MousePosition()
		shape.X, shape.Y = mx, my
		shape.Angle += seconds.FrameDelta() * 20

		var allShapes = grid.AllShapes()
		var potential = grid.ShapesAroundShape(&shape)
		for _, v := range allShapes {
			cam.DrawLinesPath(1, color.Gray, v.CornerPoints()...)
		}
		for _, v := range potential {
			cam.DrawLinesPath(2, color.Green, v.CornerPoints()...)
		}

		var col = condition.If(grid.IsCrossingShape(&shape), color.Violet, color.Cyan)
		cam.DrawLinesPath(2, col, shape.CornerPoints()...)

		var crossPoints = grid.CrossPointsWithShape(&shape)
		for _, v := range crossPoints {
			cam.DrawCircle(v[0], v[1], 3, color.Magenta)
		}
	}
}
