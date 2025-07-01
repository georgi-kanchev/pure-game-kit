package main

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/render"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = render.NewCamera()
	cam.Zoom = 3

	window.IsAntialiased = true
	var node = render.NewNode("tile", nil)
	node.AssetID = assets.LoadDefaultPatterns()

	assets.LoadDefaultSoundsUserInterface()

	for window.KeepOpen() {
		var w, h = window.Size()

		cam.SetScreenArea(0, 0, w, h)
		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(&node)
	}
}
