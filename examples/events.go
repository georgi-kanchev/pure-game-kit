package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func Events() {
	window.Create("game", true, true)
	var view = graphics.NewView(1)
	var a = ""
	for window.KeepOpen() {
		var x, y = view.PointFromEdge(0, 0)
		view.DrawText(x+10, y+10, 50, 0, palette.White, debug.MemoryUsage())
		if keyboard.IsKeyJustPressed(key.A) {
			a += "a"
		}
	}
}
