package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var mapIds = assets.LoadTiledWorld("examples/data/world.world")
	var mapGrass = tiled.NewMap(mapIds[1])

	_ = mapGrass

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DragAndZoom()
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))
	}
}

/*
func Tiled() {
	var font = assets.LoadDefaultFont()
	var cam = graphics.NewCamera(4)
	var mapIds = assets.LoadTiledWorld("examples/data/world.world")
	var grass, desert = mapIds[0], mapIds[1]
	var l1 = tilemap.LayerSprites(grass, "1", "")
	var l2 = tilemap.LayerSprites(grass, "3", "")
	var tObjs = tilemap.LayerSprites(desert, "3", "")
	var t1 = tilemap.LayerSprites(desert, "1", "")
	var g1 = tilemap.LayerShapeGrid(desert, "3", "")
	var g2 = tilemap.LayerShapeGrid(grass, "3", "")
	var pts = tilemap.LayerPoints(desert, "3", "")
	var img = tilemap.LayerSprites(desert, "7", "")
	var objs = tilemap.LayerSprites(desert, "4", "")
	var pts2 = tilemap.LayerPoints(grass, "3", "")
	var shapes = tilemap.LayerShapes(grass, "Collision", "")
	var pts3 = tilemap.LayerPoints(grass, "Collision", "")
	var imgShapes = tilemap.LayerShapes(desert, "7", "")
	var texts = tilemap.LayerTexts(grass, "Collision", "")
	var terr = tilemap.LayerSprites(mapIds[2], "Tile Layer 1", "")

	for _, text := range texts {
		text.FontId = font
	}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DragAndZoom()
		cam.DrawSprites(l1...)
		cam.DrawSprites(l2...)
		cam.DrawSprites(t1...)
		cam.DrawSprites(tObjs...)
		cam.DrawSprites(img...)
		cam.DrawSprites(objs...)
		cam.DrawSprites(terr...)
		cam.DrawTextBoxes(texts...)
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		for _, shape := range g1.All() {
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}
		for _, shape := range g2.All() {
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}

		for _, shape := range shapes {
			cam.DrawLinesPath(0.5, color.Purple, shape.CornerPoints()...)
		}
		for _, shape := range imgShapes {
			cam.DrawLinesPath(0.5, color.Magenta, shape.CornerPoints()...)
		}

		for _, pt := range pts {
			cam.DrawCircle(pt[0], pt[1], 2, color.White)
		}
		for _, pt := range pts2 {
			cam.DrawCircle(pt[0], pt[1], 2, color.Red)
		}
		for _, pt := range pts3 {
			cam.DrawCircle(pt[0], pt[1], 2, color.Purple)
		}

		if keyboard.IsKeyJustPressed(key.A) {
			tileset.TileAnimate("examples/data/objects", 198, false)
		}
		if keyboard.IsKeyJustPressed(key.S) {
			tileset.TileAnimate("examples/data/objects", 198, true)
		}

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
		}
	}
}
*/
