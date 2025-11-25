package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
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
	var desert = tiled.NewMap(mapIds[0], project)
	var sprites = desert.Sprites()

	assets.LoadDefaultFont()

	fmt.Printf("desert.Layers: %v\n", desert.Layers)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		cam.DrawSprites(sprites...)
	}
}
