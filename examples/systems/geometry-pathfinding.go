package example

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/curve"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/window"
)

func Pathfinding() {
	var view = graphics.NewView(2)
	var grid = geometry.NewShapeGrid(32, 32)

	for i := -8; i < 8; i++ {
		for j := -8; j < 8; j++ {
			if i == -1 || i == 0 || j == -1 || j == 0 {
				continue
			}

			if random.HasChance(30) {
				grid.SetAtCell(i, j, geometry.NewShapeQuad(24, 24, 0.5, 0.5))
			}
		}
	}

	var path = []float32{}

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()
		view.DrawGrid(1, 32, 32, palette.DarkGray)

		var allShapes = grid.All()
		for _, v := range allShapes {
			view.DrawLinesPath(1, palette.Gray, v.CornerPoints()...)
		}

		var mx, my = view.MousePosition()
		path = grid.FindPathDiagonally(16, 16, mx, my, false)
		path = curve.SmoothPath(path...)
		view.DrawLinesPath(1, palette.Green, path...)
		view.DrawPoints(2, palette.White, path...)
		view.DrawTextDebug(true, true, true, true)
	}
}
