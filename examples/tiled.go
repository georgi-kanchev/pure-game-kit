package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, true)
	var view = graphics.NewView(3)
	var atlasId, layerIds = assets.LoadTiledLayers("examples/data/map.tmx")

	var obj = graphics.NewTilemap(1, atlasId, layerIds[2])

	obj.Effects.BorderSize = 2
	obj.Effects.BorderColor = palette.Cyan

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		view.DrawGrid(1, 16, 16, palette.Gray)
		view.DrawObjects(&obj)
	}
}
