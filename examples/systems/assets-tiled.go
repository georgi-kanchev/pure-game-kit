package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Tiled() {
	assets.LoadFont(32, "examples/data/monogram.ttf")

	var cam = graphics.NewCamera(4)
	var mapIds = assets.LoadTiledMapsFromWorld("examples/data/world.world")
	var projectId = assets.LoadTiledProject("examples/data/game-name.tiled-project")
	var project = tiled.NewProject(projectId)
	var grass = tiled.NewMap(mapIds[0], project)
	// var sprites = desert.Sprites()

	assets.LoadDefaultFont()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()

		// grass.Draw(cam)
		cam.DrawGrid(0.5, 16, 16, color.DarkGray)

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAllTiledMaps()
			grass.Recreate()
			project.Recreate()
		}
		grass.Draw(cam)
	}
}
