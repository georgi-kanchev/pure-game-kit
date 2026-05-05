package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func MinimalRender() {
	var cam = graphics.NewCamera(1)
	for window.KeepOpen() {
		cam.DrawCircle(0, 0, 200, 32, palette.Red)
		cam.DrawTextDebug(true, true, false, true)
	}
}
