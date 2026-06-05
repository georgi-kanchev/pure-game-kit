package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, true)
	var view = graphics.NewView(1)
	var atlasId, layerIds = assets.LoadTiledLayers("examples/data/map.tmx")
	_, _ = atlasId, layerIds

	for window.KeepOpen() {
		view.DrawImage(0, 0, 1000, 1000, 0, 1, palette.White)
	}
}
