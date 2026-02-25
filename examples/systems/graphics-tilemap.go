package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/window"
)

func Tilemap() {
	var cam = graphics.NewCamera(2)
	var atlas = assets.LoadTexture("examples/data/atlas.png")
	var tilemap = graphics.NewSprite(atlas, 0, 0)

	tilemap.Effects = graphics.NewEffects()
	tilemap.Effects.TileColumns, tilemap.Effects.TileRows = 32, 32
	tilemap.Effects.TileWidth, tilemap.Effects.TileHeight = 16, 16
	tilemap.Width, tilemap.Height = 400, 400

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawSprites(tilemap)
	}
}
