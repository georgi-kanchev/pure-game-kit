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
				grid.SetAtCell(i, j, geometry.NewShapeRectangle(24, 24, 0.5, 0.5))
			}
		}
	}

	var path = [][2]float32{}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawGrid(1, 32, 32, palette.DarkGray)

		var allShapes = grid.All()
		for _, v := range allShapes {
			cam.DrawLinesPath(1, palette.Gray, v.CornerPoints()...)
		}

		var mx, my = cam.MousePosition()
		path = curve.SmoothPath(grid.FindPathDiagonally(16, 16, mx, my, false))
		cam.DrawLinesPath(1, palette.Green, path...)
		cam.DrawPoints(2, palette.White, path...)
	}
}
