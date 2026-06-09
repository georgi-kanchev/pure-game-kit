package example

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
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

	var circle = geometry.NewCircle(400, 400, 150)
	grid.AddShapes(circle)
	grid.AddShapes(geometry.NewRectangle(100, 600, 350, 200, 15))
	grid.AddShapes(geometry.NewCapsule(-200, 100, 200, 100, 60))
	grid.AddShapes(geometry.NewRoundedRectangle(700, 200, 250, 300, 45, 0.3))
	grid.AddShapes(geometry.NewLine(0, 0, 500, 500, 40))

	var x, y float32 = -500, -500

	for window.KeepOpen() {
		view.MouseDragAndZoom()

		var speed = 100 * time.Delta()
		if keyboard.IsKeyPressed(key.A) {
			x -= speed
		}
		if keyboard.IsKeyPressed(key.D) {
			x += speed
		}
		if keyboard.IsKeyPressed(key.W) {
			y -= speed
		}
		if keyboard.IsKeyPressed(key.S) {
			y += speed
		}

		if keyboard.IsKeyJustPressed(key.Enter) {
			grid.AddShapes(circle)
		}
		if keyboard.IsKeyJustPressed(key.Backspace) {
			grid.RemoveShapes(circle)
		}

		var query = geometry.NewRectangle(x, y, 150, 90, time.Running()*10)
		var neighbors = grid.Neighbors(query)

		view.DrawColor(palette.Black)
		view.DrawGrid(1, cellSize, cellSize, palette.White)

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

		for _, sh := range grid.All() {
			view.DrawShape(sh.X, sh.Y, sh.Width, sh.Height, sh.Angle, sh.Roundness, palette.LightGray)
		}
		for _, sh := range neighbors {
			view.DrawShape(sh.X, sh.Y, sh.Width, sh.Height, sh.Angle, sh.Roundness, palette.Yellow)
		}

		for _, neighbor := range neighbors {
			query = query.Collide(neighbor)
		}
		x, y = query.X, query.Y
		view.DrawShape(query.X, query.Y, query.Width, query.Height, query.Angle, query.Roundness, palette.White)
	}
}
