package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func ShapesGrids() {
	var cam = graphics.NewCamera(2)
	var grid = geometry.NewShapeGrid(32, 32)
	var shape = geometry.NewShapeCorners(
		0, 0,
		50, -20,
		100, 0,
		0, 100,
		50, 120,
		100, 100)

	for i := -8; i < 8; i++ {
		for j := -8; j < 8; j++ {
			grid.SetAtCell(i, j, geometry.NewShapeQuad(24, 24, 0.5, 0.5))
		}
	}

	for window.KeepOpen() {
		cam.MouseDragAndZoomSmoothly()
		cam.DrawGrid(1, 32, 32, palette.DarkGray)

		var mx, my = cam.MousePosition()
		shape.X, shape.Y = mx, my
		shape.Angle += time.FrameDelta() * 20

		var allShapes = grid.All()
		var potential = grid.AroundShape(shape)
		var allPts, potentailPts []float32
		for _, v := range allShapes {
			var pts = append(v.CornerPoints(), number.NaN(), number.NaN())
			allPts = append(allPts, pts...)
		}
		for _, v := range potential {
			var pts = append(v.CornerPoints(), number.NaN(), number.NaN())
			potentailPts = append(potentailPts, pts...)
		}
		cam.DrawLinesPath(1, palette.Gray, allPts...)
		cam.DrawLinesPath(2, palette.Green, potentailPts...)

		var surroundingShapes = grid.AroundShape(shape)
		var crossPoints = shape.CrossPointsWithShapes(surroundingShapes...)
		var col = condition.If(shape.IsCrossingShapes(surroundingShapes...), palette.Violet, palette.Cyan)

		cam.DrawLinesPath(2, col, shape.CornerPoints()...)

		cam.DrawPoints(3, palette.Magenta, crossPoints...)

		cam.DrawTextDebug(true, true, true, true)
	}
}
