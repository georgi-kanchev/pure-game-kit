package main

import (
	"pure-kit/engine/render"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/time"
	"pure-kit/engine/window"
)

func main() {
	var cam = render.Camera{Angle: 45, Zoom: 1}
	var angle float32 = 0.0

	window.IsAntialiased = true

	for window.KeepOpen() {
		var w, h = window.Size()

		cam.SetScreenArea(w/2, h/2, w/2, h/2)
		cam.DrawColor(color.Darken(color.Gray, 0.5))
		cam.DrawRectangle(0, 0, 200, 200, color.Red)

		var x, y = cam.CornerUpperRight(-200, 200)
		cam.DrawRectangle(x, y, 100, 100, color.Blue)

		angle += float32(time.Delta) * 10
		cam.Angle = angle
		cam.SetScreenArea(0, 0, w/2, h/2)
		cam.DrawColor(color.Darken(color.Gray, 0.75))
		cam.DrawGrid(1, 20, color.Gray)
		cam.DrawRectangle(100, 200, 200, 200, color.Green)
		cam.DrawFrame(10, color.Magenta)
	}
}
