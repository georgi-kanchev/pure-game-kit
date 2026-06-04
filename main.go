package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("split screen", false, true)

	var left, right = graphics.NewView(1), graphics.NewView(1)
	window.SetTargetFPS(60)

	var obj = graphics.NewShapeRoundedRectangle(0, 0, 400, 400, 0, 0.5)

	left.X = 400

	for window.KeepOpen() {
		var w, h = window.Size()
		left.WindowArea = graphics.NewArea(0, 0, w/2, h)
		left.Angle += time.Delta() * 30
		left.Zoom = 1.5
		left.DrawColor(palette.DarkRed)
		left.DrawObjects(&obj)
		var lx, ly = left.PointFromEdge(0, 0)
		left.DrawText(lx, ly, 100, 0, palette.White, "Hello, World!")

		var margin float32 = 300
		right.WindowArea = graphics.NewArea(w/2+margin/2, margin/2, w/2-margin, h-margin)
		right.Zoom = number.Map(number.Sine(time.Running()), -1, 1, 0.5, 4)
		right.DrawColor(palette.DarkGreen)
		right.DrawObjects(&obj)
		var rx, ry = right.PointFromEdge(0, 0)
		right.DrawText(rx, ry, 100, 0, palette.White, "Hello, World!")

		if keyboard.IsKeyJustPressed(key.F5) {
			print(debug.MemoryUsage())
		}
	}
}
