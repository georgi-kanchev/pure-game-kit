package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/window"
)

func Tiled() {
	assets.LoadFont(32, "examples/data/monogram.ttf")

	var cam = graphics.NewCamera(4)
	var mapIds = assets.LoadTiledMapsFromWorld("examples/data/world.world")
	var projectId = assets.LoadTiledProject("examples/data/game-name.tiled-project")
	var project = tiled.NewProject(projectId)
	var scene = tiled.NewScene(mapIds[0], project)

	cam.X, cam.Y = 128, 128

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()

		cam.DrawGrid(0.5, 16, 16, palette.DarkGray)
		scene.Draw(cam)

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAllTiledMaps()
			scene.Recreate()
			project.Recreate()
		}
	}
}
