package main

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = graphics.NewCamera(7)
	var parent = graphics.NewNode("")
	var tilemap = make(map[[2]float32]string, 26*21)
	var _, ids = assets.LoadDefaultAtlasRetro()
	var index = 0

	parent.X, parent.Y = -20.5, -20.5
	cam.X, cam.Y = 100, 100

	for i := range 21 {
		for j := range 26 {
			tilemap[[2]float32{float32(j), float32(i)}] = ids[index]
			index++
		}
	}
	parent.Angle = 45
	var nodemap = collection.ToPointers(graphics.NewNodesGrid(tilemap, 9, 9, &parent))
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		cam.DrawGrid(1, 9, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(nodemap...)
		cam.DrawNodes(&parent)

		var child = nodemap[5]
		child.X, child.Y = parent.MousePosition(&cam)
	}
}
