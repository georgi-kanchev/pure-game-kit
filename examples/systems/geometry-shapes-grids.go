package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func ShapesGrids() {
	var cam = graphics.NewCamera(2)
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
			grid.SetAtCell(i, j, geometry.NewShapeRectangle(24, 24, 0.5, 0.5))
		}
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawGrid(1, 32, 32, palette.DarkGray)

		var mx, my = cam.MousePosition()
		shape.X, shape.Y = mx, my
		shape.Angle += time.FrameDelta() * 20

		var allShapes = grid.All()
		var potential = grid.AroundShape(shape)
		for _, v := range allShapes {
			cam.DrawLinesPath(1, palette.Gray, v.CornerPoints()...)
		}
		for _, v := range potential {
			cam.DrawLinesPath(2, palette.Green, v.CornerPoints()...)
		}

		var surroundingShapes = grid.AroundShape(shape)
		var crossPoints = shape.CrossPointsWithShapes(surroundingShapes...)
		var col = condition.If(shape.IsCrossingShapes(surroundingShapes...), palette.Violet, palette.Cyan)

		cam.DrawLinesPath(2, col, shape.CornerPoints()...)

		for _, v := range crossPoints {
			cam.DrawCircle(v[0], v[1], 3, palette.Magenta)
		}
	}
}
