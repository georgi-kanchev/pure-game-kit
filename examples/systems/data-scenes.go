package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/tiled/tilemap"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var layer1, layer2, objs, t1 = reload()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DragAndZoom()
		cam.DrawSprites(layer1...)
		cam.DrawSprites(layer2...)
		cam.DrawSprites(t1...)
		cam.DrawSprites(objs...)
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		if keyboard.IsKeyPressedOnce(key.F5) {
			layer1, layer2, objs, t1 = reload()
		}
	}
}

func reload() (layer1, layer2, objs, t1 []*graphics.Sprite) {
	var mapIds = assets.LoadTiledWorlds("examples/data/world.world")
	assets.LoadTiledTilesets("examples/data/atlas.tsx")
	assets.LoadTiledTilesets("examples/data/objects.tsx")
	layer1 = tilemap.LayerTiles(mapIds[0], "Tile Layer 1")
	layer2 = tilemap.LayerTiles(mapIds[0], "Tile Layer 2")
	objs = tilemap.LayerTiles(mapIds[1], "Objects")
	t1 = tilemap.LayerTiles(mapIds[1], "Tile Layer 1")
	return
}
