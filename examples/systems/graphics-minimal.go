package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func MinimalRender() {
	var view = graphics.NewView(1)
	for window.KeepOpen() {
		view.DrawCircle(0, 0, 200, 32, palette.Red)
		view.DrawTextDebug(true, true, false, true)
	}
}
