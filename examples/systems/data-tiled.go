package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/tiled/tilemap"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var layer1, layer2, objs, t1, g1, g2 = reload()
	var grid1, grid2 = g1.All(), g2.All()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DragAndZoom()
		cam.DrawSprites(layer1...)
		cam.DrawSprites(layer2...)
		cam.DrawSprites(t1...)
		cam.DrawSprites(objs...)
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		for _, shape := range grid1 {
			var cellX, cellY = g1.Cell(shape)
			_, _ = cellX, cellY
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}
		for _, shape := range grid2 {
			var cellX, cellY = g2.Cell(shape)
			_, _ = cellX, cellY
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}

		if keyboard.IsKeyPressedOnce(key.F5) {
			layer1, layer2, objs, t1, g1, g2 = reload()
			grid1, grid2 = g1.All(), g2.All()
		}
	}
}

func reload() (layer1, layer2, objs, t1 []*graphics.Sprite, g1, g2 *geometry.ShapeGrid) {
	var mapIds = assets.LoadTiledWorld("examples/data/world.world")
	assets.LoadTiledTileset("examples/data/atlas.tsx")
	assets.LoadTiledTileset("examples/data/objects.tsx")
	layer1 = tilemap.LayerTiles(mapIds[0], "1")
	layer2 = tilemap.LayerTiles(mapIds[0], "3")
	objs = tilemap.LayerTiles(mapIds[1], "3")
	t1 = tilemap.LayerTiles(mapIds[1], "1")
	g1 = tilemap.LayerTilesShapeGrid("examples/data/desert", "3", "")
	g2 = tilemap.LayerTilesShapeGrid("examples/data/map", "3", "")
	return
}
