package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, true)
	var view = graphics.NewView(1)
	var atlasId, layerIds = assets.LoadTiledLayers("examples/data/map.tmx")

	var obj = graphics.NewTilemap(3, atlasId, layerIds[3])

	for window.KeepOpen() {
		view.DrawObjects(&obj)
	}
}
