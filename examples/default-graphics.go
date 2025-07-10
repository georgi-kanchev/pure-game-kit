package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/window"
)

func DefaultGraphics() {
	var camera = graphics.NewCamera(7)
	var _, _ = assets.LoadDefaultAtlasUI()
	var node = graphics.NewNode("!")

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.DrawNodes(&node)
	}
}
