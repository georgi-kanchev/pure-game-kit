package example

import (
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/motion/curve"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/random"
	"pure-game-kit/window"
)

func Pathfinding() {
	var cam = graphics.NewCamera(2)
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
		cam.MouseDragAndZoomSmoothly()
		cam.DrawGrid(1, 32, 32, palette.DarkGray)

		var allShapes = grid.All()
		for _, v := range allShapes {
			cam.DrawLinesPath(1, palette.Gray, v.CornerPoints()...)
		}

		var mx, my = cam.MousePosition()
		path = grid.FindPathDiagonally(16, 16, mx, my, false)
		path = curve.SmoothPath(path...)
		cam.DrawLinesPath(1, palette.Green, path...)
		cam.DrawPoints(2, palette.White, path...)
		cam.DrawTextDebug(true, true, true, true)
	}
}
