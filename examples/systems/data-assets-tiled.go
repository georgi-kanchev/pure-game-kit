package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var tileSetId string
	var tileDataIds []string
	var pts [][2]float32
	var hotreload = func() {
		tileSetId = assets.LoadTileSet("examples/data/atlas.png", 16, 16)
		tileDataIds = assets.LoadTiledData("examples/data/map.tmx")
		pts = assets.LoadTiledPoints("examples/data/map.tmx", "Objects")
	}

	hotreload()

	cam.X, cam.Y = 128, 128

	var tileMaps = make([]*graphics.TileMap, len(tileDataIds))
	for i, t := range tileDataIds {
		tileMaps[i] = graphics.NewTileMap(tileSetId, t)
		tileMaps[i].PivotX, tileMaps[i].PivotY = 0, 0
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tileMaps...)
		cam.DrawShapesFast(palette.Red, pts...)
		cam.DrawPoints(2, palette.White, pts...)

		if keyboard.IsKeyJustPressed(key.F5) {
			hotreload()
		}
	}
}
