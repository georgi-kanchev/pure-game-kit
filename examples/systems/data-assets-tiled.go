package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var tileSetId string
	var tileDataIds []string
	var hotreload = func() {
		tileSetId, tileDataIds = assets.LoadTiledLayers("examples/data/map.tmx")
	}

	hotreload()

	cam.X, cam.Y = 128, 128

	var tileMaps = make([]*graphics.TileMap, len(tileDataIds))
	for i, t := range tileDataIds {
		tileMaps[i] = graphics.NewTileMap(tileSetId, t)
		tileMaps[i].PivotX, tileMaps[i].PivotY = 0, 0
		tileMaps[i].Angle = 5
	}

	for window.KeepOpen() {
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tileMaps...)

		cam.DrawShapes(color.FadeOut(palette.Red, 0.5), tileMaps[1].Points()...)
		cam.DrawShapes(color.FadeOut(palette.Blue, 0.5), tileMaps[3].Points()...)

		if mouse.IsButtonPressed(button.Left) {
			var x, y = tileMaps[3].CellAtPoint(cam.MousePosition())
			tileMaps[3].SetTile(x, y, graphics.NewTile(22))
		}

		if keyboard.IsKeyJustPressed(key.F5) {
			hotreload()
		}
		cam.DrawTextDebug(true, true, true, true)
	}
}
