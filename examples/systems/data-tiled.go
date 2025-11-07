package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled/tilemap"
	"pure-game-kit/tiled/tileset"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

func Tiled() {
	var font = assets.LoadDefaultFont()
	var cam = graphics.NewCamera(4)
	var l1, l2, tObjs, t1, img, objs, terr []*graphics.Sprite
	var texts []*graphics.TextBox
	var pts, pts2, pts3 [][2]float32
	var g1, g2 *geometry.ShapeGrid
	var shapes, imgShapes []*geometry.Shape
	var reload = func() {
		var mapIds = assets.LoadTiledWorld("examples/data/world.world")
		var grass, desert = mapIds[0], mapIds[1]
		l1 = tilemap.LayerSprites(grass, "1", "")
		l2 = tilemap.LayerSprites(grass, "3", "")
		tObjs = tilemap.LayerSprites(desert, "3", "")
		t1 = tilemap.LayerSprites(desert, "1", "")
		g1 = tilemap.LayerShapeGrid(desert, "3", "")
		g2 = tilemap.LayerShapeGrid(grass, "3", "")
		pts = tilemap.LayerPoints(desert, "3", "")
		img = tilemap.LayerSprites(desert, "7", "")
		objs = tilemap.LayerSprites(desert, "4", "")
		pts2 = tilemap.LayerPoints(grass, "3", "")
		shapes = tilemap.LayerShapes(grass, "Collision", "")
		pts3 = tilemap.LayerPoints(grass, "Collision", "")
		imgShapes = tilemap.LayerShapes(desert, "7", "")
		texts = tilemap.LayerTexts(grass, "Collision", "")
		terr = tilemap.LayerSprites(mapIds[2], "Tile Layer 1", "")

		for _, text := range texts {
			text.FontId = font
		}
	}
	reload()

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

		if keyboard.IsKeyJustPressed(key.Y) {
			assets.UnloadTexture("examples/data/objects")
		}
		if keyboard.IsKeyJustPressed(key.T) {
			assets.LoadTexture("examples/data/objects.png")
		}

		if keyboard.IsKeyJustPressed(key.F5) {
			reload()
		}
	}
}
