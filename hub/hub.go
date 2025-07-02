package main

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = graphics.NewCamera(3)
	var node = graphics.NewNode("")

	assets.LoadDefaultTexture()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&node)
	}
}
