package main

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/render"
	"pure-kit/engine/tiles"
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func main() {
	var cam = render.NewCamera()
	cam.Zoom = 3

	window.IsAntialiased = true
	var node = render.NewNode("tile", nil)
	node.AssetID = "default-atlas"

	assets.LoadDefaults(true, true)
	assets.LoadTexturesFromFiles("hell.png")
	assets.LoadAtlasFromTexture("hell", 32, 32, 0)

	assets.LoadTileFromAtlas("hell", "tile", 4, 1, 1, 1)
	assets.LoadTileFromAtlas("hell", "tile2", 0, 0, 1, 1)

	var tilemap = tiles.Map{}
	tilemap.SetTile(0, 0, "tile")
	tilemap.SetTile(1, 0, "tile")
	tilemap.SetTile(-2, 0, "tile")
	tilemap.SetTile(2, 0, "tile2")
	tilemap.SetTile(0, 2, "tile2")
	var tilemapRender = collection.ToPointers(render.NewNodesGrid(tilemap.Tiles, 32, 32, nil))

	for window.KeepOpen() {
		var w, h = window.Size()

		cam.SetScreenArea(0, 0, w, h)
		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))
		cam.DrawNodes(tilemapRender...)
		cam.DrawNodes(&node)
	}
}
