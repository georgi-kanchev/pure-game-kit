package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var tileSetId string
	var tileDataIds []string
	var pts []float32
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
		cam.DrawShapes(palette.Red, pts...)
		cam.DrawPoints(2, palette.White, pts...)
		cam.DrawTextDebug(true, true, true, true)

		cam.DrawShapes(palette.Red, tileMaps[3].Points()...)

		if mouse.IsButtonPressed(button.Left) {
			var x, y = tileMaps[3].CellAtPoint(cam.MousePosition())
			tileMaps[3].SetTile(x, y, graphics.NewTile(22))
		}
		// fmt.Printf("cell: %v, %v\n", x, y)

		// var tilePts = tileMaps[3].PointsAtCell(x, y)
		// cam.DrawShapes(palette.Blue, tilePts...)

		if keyboard.IsKeyJustPressed(key.F5) {
			hotreload()
		}
	}
}
