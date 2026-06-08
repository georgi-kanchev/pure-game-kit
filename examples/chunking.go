package example

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func Chunking() {
	window.Create("example - shape grid chunking", true, true)

	var view = graphics.NewView(1)

	const cellSize float32 = 100
	var grid = geometry.NewShapeGrid(cellSize)

	grid.SetAtCell(0, 0, geometry.NewCircle(50, 50, 40))
	grid.SetAtCell(2, 0, geometry.NewRectangle(250, 50, 80, 50, 30))
	grid.SetAtCell(-1, 1, geometry.NewCapsule(-50, 150, -100, 150, 15))
	grid.SetAtCell(3, 1, geometry.NewRoundedRectangle(350, 150, 90, 60, 0, 0.4))
	grid.SetAtCell(0, 2, geometry.NewLine(80, 280, 30, 250, 12))
	grid.SetAtCell(2, 2, geometry.NewCircle(250, 250, 35))
	grid.SetAtCell(4, -1, geometry.NewRectangle(450, -50, 50, 70, -15))

	for window.KeepOpen() {
		view.MouseDragAndZoom()

		var mx, my = view.MousePosition()
		var query = geometry.NewRectangle(mx, my, 150, 90, time.Running()*10)
		var neighbors = grid.Neighbors(query)

		view.DrawColor(palette.Black)
		view.DrawGrid(1, cellSize, cellSize, palette.DarkGray)

		var qMinX, qMinY, qW, qH = query.Bounds()
		var qMaxX, qMaxY = qMinX + qW, qMinY + qH
		var sx = int(number.RoundDown(qMinX / cellSize))
		var sy = int(number.RoundDown(qMinY / cellSize))
		var ex = int(number.RoundDown(qMaxX / cellSize))
		var ey = int(number.RoundDown(qMaxY / cellSize))
		for cx := sx; cx <= ex; cx++ {
			for cy := sy; cy <= ey; cy++ {
				view.DrawShape(float32(cx)*cellSize+cellSize/2, float32(cy)*cellSize+cellSize/2, cellSize, cellSize, 0, 0, palette.Gray)
			}
		}

		var isNeighbor = make(map[geometry.Shape]bool, len(neighbors))
		for _, n := range neighbors {
			isNeighbor[n] = true
		}

		for _, sh := range grid.All() {
			var color = palette.LightGray
			if isNeighbor[sh] {
				color = palette.Yellow
			}
			view.DrawShape(sh.X, sh.Y, sh.Width, sh.Height, sh.Angle, sh.Roundness, color)
		}

		view.DrawShape(query.X, query.Y, query.Width, query.Height, query.Angle, query.Roundness, palette.White)
	}
}
