package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/random"
	"pure-game-kit/window"
)

func Tilemaps() {
	var cam = graphics.NewCamera(2)
	var atlasId = assets.LoadTileSet("examples/data/atlas.png", 16, 16)
	var tileDataId = assets.LoadTileData("tilemap", 320, 320)
	var tilemap = graphics.NewTileMap(atlasId, tileDataId)

	// tilemap.SetTileArea(0, 0, 320, 320, graphics.NewTile(29))

	// tilemap.Effects = graphics.NewEffects()
	// tilemap.Effects.Saturation = 0.8

	for y := range 320 {
		for x := range 320 {
			tilemap.SetTile(x, y, graphics.NewTile(random.Range[uint16](0, 335)))
		}
	}

	for window.KeepOpen() {
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tilemap)
		cam.DrawTextDebug(true, true, true, true)

		if mouse.IsButtonPressed(button.Left) {
			var mx, my = cam.MousePosition()
			var x, y = tilemap.PointToLocal(mx, my)
			tilemap.SetTile(int(x/16), int(y/16), graphics.NewTileAnimated(106, 15, byte(x/16), 20))
		}
	}
}
