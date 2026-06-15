package example

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func CollisionGrid() {
	window.Create("example - collision grid", false, true)

	var view = graphics.NewView(1)

	const cellSize float32 = 50
	var grid = geometry.NewShapeGrid(cellSize)

	var circle = geometry.NewCircle(400, 400, 150)
	grid.AddShapes(circle)
	grid.AddShapes(geometry.NewRectangle(100, 600, 350, 200, 15))
	grid.AddShapes(geometry.NewCapsule(-200, 100, 200, 100, 60))
	grid.AddShapes(geometry.NewRoundedRectangle(700, 200, 250, 300, 45, 0.3))
	grid.AddShapes(geometry.NewLine(0, 0, 500, 500, 40))

	var player = geometry.NewRectangle(-300, -300, 150, 90, 0)

	window.SetTargetFPS(0)

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		var speed = 300 * time.Delta()
		if keyboard.IsKeyPressed(key.A) {
			player.X -= speed
		}
		if keyboard.IsKeyPressed(key.D) {
			player.X += speed
		}
		if keyboard.IsKeyPressed(key.W) {
			player.Y -= speed
		}
		if keyboard.IsKeyPressed(key.S) {
			player.Y += speed
		}
		player.Angle = time.Running() * 20

		if keyboard.IsKeyJustPressed(key.Enter) {
			grid.AddShapes(circle)
		}
		if keyboard.IsKeyJustPressed(key.Backspace) {
			grid.RemoveShapes(circle)
		}

		var neighbors = grid.Neighbors(player)
		for _, neighbor := range neighbors {
			player = player.Collide(neighbor)
		}

		for y := -10; y < 20; y++ {
			for x := -10; x < 20; x++ {
				var sh = grid.AtCell(x, y)
				var col = palette.DarkRed
				if len(sh) > 0 {
					col = palette.DarkGreen
				}
				view.DrawShape(float32(x)*cellSize+cellSize/2, float32(y)*cellSize+cellSize/2, cellSize, cellSize, 0, 0, col, graphics.Area{})
			}
		}

		var bx, by, bw, bh = player.Bounds()
		view.DrawShape(bx+bw/2, by+bh/2, bw, bh, 0, 0, color.RGBA(0, 0, 0, 128), graphics.Area{})

		view.DrawGrid(1, cellSize, cellSize, palette.Gray)

		for _, sh := range grid.All() {
			view.DrawShape(sh.X, sh.Y, sh.Width, sh.Height, sh.Angle, sh.Roundness, palette.LightGray, graphics.Area{})
		}
		for _, sh := range neighbors {
			view.DrawShape(sh.X, sh.Y, sh.Width, sh.Height, sh.Angle, sh.Roundness, palette.Yellow, graphics.Area{})
		}
		view.DrawShape(player.X, player.Y, player.Width, player.Height, player.Angle, player.Roundness, palette.White, graphics.Area{})

		var mx, my = view.MousePosition()
		var pts = grid.FindPathDiagonally(player.X, player.Y, mx, my, true)
		view.DrawPath(pts, 5, palette.Cyan, graphics.Area{})
		view.DrawDebugInfo(true)
	}
}
