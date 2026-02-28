package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Tilemap() {
	var cam = graphics.NewCamera(2)
	var atlasId = assets.LoadTileAtlas("examples/data/atlas.png", 16, 16)
	var tileDataId = assets.LoadTileData("tilemap", 2048, 2048)
	var tilemap = graphics.NewTileMap(atlasId, tileDataId)

	window.FrameRateLimit = 0

	assets.SetTileArea(tileDataId, 0, 0, 2048, 2048, 29)

	var fps = ""

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tilemap)

		if mouse.IsButtonPressed(button.Left) {
			var mx, my = cam.MousePosition()
			var x, y = tilemap.PointToLocal(mx, my)
			assets.SetTile(tileDataId, int(x/16), int(y/16), 106)
		}

		var tlx, tly = cam.PointFromEdge(0, 0)
		cam.DrawText(fps, tlx, tly, 150/cam.Zoom)

		if condition.TrueEvery(0.1, "fps") {
			fps = text.New("Current FPS: ", time.FrameRate(), "\n", "Average FPS: ", time.FrameRateAverage())
		}
	}
}
