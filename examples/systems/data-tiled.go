package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/tiled/tilemap"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/time"
	"pure-kit/engine/utility/time/unit"
	"pure-kit/engine/window"
)

func Tiled() {
	var cam = graphics.NewCamera(4)
	var l1, l2, tObjs, t1, img, objs []*graphics.Sprite
	var pts, pts2 [][2]float32
	var g1, g2 *geometry.ShapeGrid
	var reload = func() {
		var mapIds = assets.LoadTiledWorld("examples/data/world.world")
		var grass, desert = mapIds[0], mapIds[1]
		assets.LoadTiledTileset("examples/data/atlas.tsx")
		assets.LoadTiledTileset("examples/data/objects.tsx")
		l1 = tilemap.LayerSprites(grass, "1")
		l2 = tilemap.LayerSprites(grass, "3")
		tObjs = tilemap.LayerSprites(desert, "3")
		t1 = tilemap.LayerSprites(desert, "1")
		g1 = tilemap.LayerShapeGrid(desert, "3", "")
		g2 = tilemap.LayerShapeGrid(grass, "3", "")
		pts = tilemap.LayerPoints(desert, "3", "")
		img = tilemap.LayerSprites(desert, "7")
		objs = tilemap.LayerSprites(desert, "4")
		pts2 = tilemap.LayerPoints(grass, "3", "")
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
		cam.DrawGrid(0.5, 16, 16, color.Darken(color.Gray, 0.5))

		for _, shape := range g1.All() {
			var cellX, cellY = g1.Cell(shape)
			_, _ = cellX, cellY
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}
		for _, shape := range g2.All() {
			var cellX, cellY = g2.Cell(shape)
			_, _ = cellX, cellY
			cam.DrawLinesPath(0.5, color.Red, shape.CornerPoints()...)
		}

		for _, pt := range pts {
			cam.DrawCircle(pt[0], pt[1], 2, color.White)
		}
		for _, pt := range pts2 {
			cam.DrawCircle(pt[0], pt[1], 2, color.Red)
		}

		if keyboard.IsKeyPressedOnce(key.A) {
			fmt.Printf("%v\n", time.AsClock24(187368730, ":", unit.All))
		}

		if keyboard.IsKeyPressedOnce(key.F5) {
			reload()
		}
	}
}
