package main

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = graphics.NewCamera(7)
	var parent = graphics.NewNode("")

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 9, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&parent)
	}
}
