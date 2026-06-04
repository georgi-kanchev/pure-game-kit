package main

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("split screen", false, true)

	var viewLeft, viewRight = graphics.NewView(1), graphics.NewView(1)
	window.SetTargetFPS(60)

	for window.KeepOpen() {
		var w, h = window.Size()
		viewLeft.WindowArea = graphics.NewArea(0, 0, w/2, h)
		viewLeft.DrawColor(palette.DarkRed)
		viewLeft.DrawText(0, 0, 100, 0, palette.Red, "Hello, World!")

		viewRight.WindowArea = graphics.NewArea(w/2+50, 50, w/2-100, h-100)
		viewRight.DrawColor(palette.DarkGreen)
		var x, y = viewRight.MousePosition()
		viewRight.DrawText(x, y, 100, 0, palette.Red, text.New("x: ", x, "y: ", y))

		if keyboard.IsKeyJustPressed(key.F5) {
			print(debug.MemoryUsage())
		}
	}
}
