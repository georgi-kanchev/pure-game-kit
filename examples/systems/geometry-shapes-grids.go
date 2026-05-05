package example

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func ShapesGrids() {
	var view = graphics.NewView(2)
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
		view.MouseDragAndZoomSmoothly()
		view.DrawGrid(1, 32, 32, palette.DarkGray)

		var mx, my = view.MousePosition()
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
		view.DrawLinesPath(1, palette.Gray, allPts...)
		view.DrawLinesPath(2, palette.Green, potentailPts...)

		var surroundingShapes = grid.AroundShape(shape)
		var crossPoints = shape.CrossPointsWithShapes(surroundingShapes...)
		var col = condition.If(shape.IsCrossingShapes(surroundingShapes...), palette.Violet, palette.Cyan)

		view.DrawLinesPath(2, col, shape.CornerPoints()...)

		view.DrawPoints(3, palette.Magenta, crossPoints...)

		view.DrawTextDebug(true, true, true, true)
	}
}
