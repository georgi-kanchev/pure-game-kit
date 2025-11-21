package example

import (
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
	var desert = tiled.NewMap(mapIds[1], project)

	var sprites = desert.Layers[0].Sprites()
	var shapes = desert.Layers[0].Shapes()
	var points = desert.Layers[0].Points()
	var lines = desert.Layers[0].Lines()

	assets.LoadDefaultFont()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmooth()
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		cam.DrawSprites(sprites...)
		cam.DrawPoints(1, color.Red, points...)
		cam.DrawLinesPath(1, color.Blue, lines...)

		for _, shape := range shapes {
			cam.DrawShapes(color.FadeOut(color.Green, 0.75), shape.CornerPoints()...)
		}
	}
}
