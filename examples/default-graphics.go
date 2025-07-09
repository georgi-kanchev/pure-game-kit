package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func DefaultGraphics() {
	var camera = graphics.NewCamera(7)
	var retro, _ = assets.LoadDefaultAtlasRetro()
	var retroNode = graphics.NewNode(retro)

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.ScreenWidth /= 2
		camera.DrawGrid(1, 9, color.Darken(color.Gray, 0.5))
		camera.DrawNodes(&retroNode)
		camera.DragAndZoom()
	}
}
