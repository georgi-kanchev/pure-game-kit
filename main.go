package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("split screen", false, true)

	var left, right = graphics.NewView(1), graphics.NewView(1)
	window.SetTargetFPS(60)

	var obj = graphics.NewShapeRoundedRectangle(0, 0, 400, 400, 0, 0.5)

	left.X = 0

	for window.KeepOpen() {
		// obj.Angle += time.Delta() * 30

		var w, h = window.Size()
		left.WindowArea = graphics.NewArea(0, 0, w/2, h)
		left.Angle += time.Delta() * 30
		left.Zoom = 1.5
		left.DrawColor(palette.DarkRed)
		left.DrawObjects(&obj)

		var x, y = left.PointFromEdge(0, 0)
		left.DrawText(x, y, 50, 0, palette.White, "Hello, World!")

		right.WindowArea = graphics.NewArea(w/2+50, 50, w/2-100, h-100)
		right.DrawColor(palette.DarkGreen)
		right.DrawObjects(&obj)
		right.DrawText(0, 0, 50, 0, palette.White, "Hello, World!")

		if keyboard.IsKeyJustPressed(key.F5) {
			print(debug.MemoryUsage())
		}
	}
}
