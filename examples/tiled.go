package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, false)
	var view = graphics.NewView(3)
	var layerIds = assets.LoadTileLayersFromTiled("examples/data/map.tmx")

	var layers []graphics.Object
	for _, id := range layerIds {
		var layer = graphics.NewTilemap(1, id)
		layers = append(layers, layer)
	}

	var shapes, cellShapes = collection.NewList[geometry.Shape](), collection.NewList[geometry.Shape]()
	layers[1].TilemapShapes(shapes.AsSlice())
	layers[3].TilemapShapes(cellShapes.AsSlice())

	layers[0].TileLayerId.SetTile(0, 0, assets.NewTile(55))

	window.SetTargetFPS(0)

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		for _, l := range layers {
			view.DrawObject(&l)
		}
		for _, s := range *shapes.AsSlice() {
			view.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, palette.Red, geometry.Area{})
		}
		for _, s := range *cellShapes.AsSlice() {
			view.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, palette.DarkRed, geometry.Area{})
		}

		// view.DrawGrid(0.3, 16, 16, palette.DarkGray)
		view.DrawDebugInfo(true)
	}
}
