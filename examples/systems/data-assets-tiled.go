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
	var pts []float32
	var hotreload = func() {
		tileSetId, tileDataIds = assets.LoadTiledData("examples/data/map.tmx")
		pts = assets.LoadTiledPoints("examples/data/map.tmx", "Objects", "Tile Layer 3")
	}

	hotreload()

	cam.X, cam.Y = 128, 128

	var tileMaps = make([]*graphics.TileMap, len(tileDataIds))
	for i, t := range tileDataIds {
		tileMaps[i] = graphics.NewTileMap(tileSetId, t)
		tileMaps[i].PivotX, tileMaps[i].PivotY = 0, 0
	}

	for window.KeepOpen() {
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tileMaps...)
		cam.DrawShapes(palette.Red, pts...)
		cam.DrawPoints(2, palette.White, pts...)
		cam.DrawTextDebug(true, true, true, true)

		if keyboard.IsKeyJustPressed(key.F5) {
			hotreload()
		}
	}
}
