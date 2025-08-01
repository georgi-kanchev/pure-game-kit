package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func MinimalRender() {
	var cam = graphics.NewCamera(1)
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawCircle(0, 0, 100, color.Red)
	}
}
