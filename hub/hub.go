package main

import (
	texture "pure-kit/engine/data/assets"
	"pure-kit/engine/render"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = render.NewCamera()
	cam.Zoom = 8

	window.IsAntialiased = true

	var node = render.NewNode("flipped", nil)
	node.OriginX, node.OriginY = 0, 0
	texture.LoadTexturesFromFiles("hell.png")
	texture.LoadAtlasFromTexture("hell", 32, 32, 0)

	texture.LoadCellFromAtlas("hell", "flipped", 6, 1, -1, 1)

	for window.KeepOpen() {
		var w, h = window.Size()

		cam.SetScreenArea(0, 0, w, h)
		cam.DrawColor(color.Darken(color.Gray, 0.75))
		cam.DrawGrid(1, 32, color.Gray)
		cam.DrawNode(&node)
	}
}
