package example

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, true)
	var view = graphics.NewView(3)
	var layerIds = assets.LoadTileLayersFromTiled("examples/data/map.tmx")

	var layers []graphics.Object
	for _, id := range layerIds {
		var layer = graphics.NewTilemap(1, id)
		layers = append(layers, layer)
	}

	var shapes = layers[1].TilemapShapes()
	var cellShapes = layers[3].TilemapShapes()

	layers[0].TileLayerId.SetTile(0, 0, assets.NewTile(55))

	fmt.Printf("debug.LinesOfCode(): %v\n", debug.LinesOfCode())

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		for _, l := range layers {
			view.DrawObject(&l)
		}
		for _, s := range shapes {
			view.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, palette.Red, graphics.Area{})
		}
		for _, s := range cellShapes {
			view.DrawShape(s.X, s.Y, s.Width, s.Height, s.Angle, s.Roundness, palette.DarkRed, graphics.Area{})
		}

		view.DrawGrid(0.3, 16, 16, palette.DarkGray)
		view.DrawDebugInfo(true)
	}
}
