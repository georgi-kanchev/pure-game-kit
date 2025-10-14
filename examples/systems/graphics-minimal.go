package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Minimal() {
	var cam = graphics.NewCamera(1)
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawColor(color.Darken(color.Gray, 0.5))
		cam.DrawCircle(0, 0, 100, color.Red)

		var x, y = cam.MousePosition()
		cam.DrawRectangle(x, y, 200, 200, 10, color.Blue, color.Red, color.Yellow, color.White)
	}
}
