package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Tilemap() {
	var cam = graphics.NewCamera(2)
	var atlasId = assets.LoadTileSet("examples/data/atlas.png", 16, 16)
	var tileDataId = assets.LoadTileData("tilemap", 32, 32)
	var tilemap = graphics.NewTileMap(atlasId, tileDataId)

	tilemap.SetTileArea(0, 0, 32, 32, graphics.NewTile(29))

	var fps = ""

	tilemap.Effects = graphics.NewEffects()
	tilemap.Effects.Saturation = 0.8

	window.FrameRateLimit = 0

	for y := range 32 {
		for x := range 32 {
			tilemap.SetTile(x, y, graphics.NewTile(random.Range[uint16](0, 335)))
		}
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tilemap)

		if mouse.IsButtonPressed(button.Left) {
			var mx, my = cam.MousePosition()
			var x, y = tilemap.PointToLocal(mx, my)
			tilemap.SetTile(int(x/16), int(y/16), graphics.NewTileAnimated(106, 15, byte(x/16), 20))
			var tile = tilemap.TileAt(int(x/16), int(y/16))
			debug.Print(text.New(tile))
		}

		if condition.TrueEvery(0.1, "fps") {
			fps = text.New("Current FPS: ", time.FrameRate(), "\n", "Average FPS: ", time.FrameRateAverage())
		}

		var tlx, tly = cam.PointFromEdge(0, 0)
		cam.DrawText(fps, tlx, tly, 50/cam.Zoom)
	}
}
