package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/tiled"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Scenes() {
	var cam = graphics.NewCamera(7)
	var sprite = graphics.NewSprite("examples/data/atlas[335]", 0, 0)
	var mapData = assets.LoadTiledMaps("examples/data/map.tmx")[0]
	var props = tiled.ExtractMapProperties(mapData)
	var world = assets.LoadTiledWorlds("examples/data/world.world")[0]
	var tileset = assets.LoadTiledTilesets("examples/data/atlas.tsx")[0]

	fmt.Printf("props: %v\n", props)
	fmt.Printf("worlds: %v\n", world)
	fmt.Printf("tilesets: %v\n", tileset)

	var myNumber = tiled.TilesetProperty(tileset, tiled.TilesetColumns)
	fmt.Printf("myNumber: %v\n", myNumber)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 9, 9, color.Darken(color.Gray, 0.5))
		cam.DrawSprites(&sprite)
	}
}
