package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/window"
)

func Tilemap() {
	var cam = graphics.NewCamera(2)
	var atlasId = assets.LoadTileAtlas("examples/data/atlas.png", 16, 16)
	var tileDataId = assets.LoadTileData("tilemap", 32, 32)
	var tilemap = graphics.NewTileMap(atlasId, tileDataId)

	assets.SetTileArea(tileDataId, 0, 0, 32, 32, 29, 0, false)

	var ang = 0
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()
		cam.DrawTileMaps(tilemap)

		if mouse.IsButtonJustPressed(button.Left) {
			var mx, my = cam.MousePosition()
			var x, y = tilemap.PointToLocal(mx, my)
			var tile, rot, flip = assets.Tile(tileDataId, int(x/16), int(y/16))
			fmt.Printf("%v %v %v\n", tile, rot, flip)
			assets.SetTile(tileDataId, int(x/16), int(y/16), 106, ang, true)
			ang++
		}
	}
}
