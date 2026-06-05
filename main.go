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
		obj.Roundness = number.Map(number.Sine(time.Running()/3), -1, 1, 0, 1)

		var w, h = window.Size()
		left.WindowArea = graphics.NewArea(0, 0, w/2, h)
		left.Angle += time.Delta() * 10
		left.Zoom = 1.5
		left.DrawColor(palette.DarkGray)
		left.DrawGrid(2, 100, 100, palette.Gray)
		left.DrawObjects(&obj)
		var lx, ly = left.PointFromScreen(left.WindowArea.X+10, left.WindowArea.Y+10)
		left.DrawText(lx, ly, 100, 0, palette.White, "Left View")

		var margin float32 = 300
		right.WindowArea = graphics.NewArea(w/2-margin/2, margin/2, w/2-margin, h-margin)
		right.Zoom = number.Map(number.Sine(time.Running()/3), -1, 1, 0.5, 4)
		right.DrawColor(palette.Gray)
		right.DrawGrid(2, 100, 100, palette.LightGray)
		right.DrawObjects(&obj)
		var rx, ry = right.PointFromScreen(right.WindowArea.X+10, right.WindowArea.Y+10)
		right.DrawText(rx, ry, 100, 0, palette.White, "Right View")

		if keyboard.IsKeyJustPressed(key.F5) {
			print(debug.MemoryUsage())
		}
	}
}
