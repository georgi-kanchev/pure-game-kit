package example

import (
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func MinimalRender() {
	var cam = graphics.NewCamera(1)
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawCircle(0, 0, 100, color.Red)
	}
}
