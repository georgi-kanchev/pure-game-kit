package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/window"
)

func Collisions() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawNodes(&node)
	}
}
