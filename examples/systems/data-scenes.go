package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/geometry"
	"pure-kit/engine/graphics"
	"pure-kit/engine/tiled/tileset"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Scenes() {
	var cam = graphics.NewCamera(20)
	var id = assets.LoadTiledTilesets("examples/data/atlas.tsx")[0]
	var sprite = graphics.NewSprite(id+"[115]", 0, 0)

	var myNumber = tileset.Property(id, tileset.PropertyColumns)
	fmt.Printf("myNumber: %v\n", myNumber)

	var points = tileset.TileShapeCorners(id, 115, "collision")
	var x, y = tileset.TileShapePoint(id, 115, "5")
	var solid = tileset.TileShapeProperty(id, 115, "collision", "solid")

	fmt.Printf("solid: %v\n", solid)

	fmt.Printf("points: %v\n", points)

	fmt.Printf("%v, %v\n", x, y)

	var sh = geometry.NewShapeCorners(points...)
	sprite.Width, sprite.Height = 16, 16
	sprite.PivotX, sprite.PivotY = 0, 0

	var durs = tileset.TileAnimationTileIds(id, 139)

	fmt.Printf("durs: %v\n", durs)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 9, 9, color.Darken(color.Gray, 0.5))
		cam.DrawSprites(&sprite)
		cam.DrawLinesPath(0.2, color.White, sh.CornerPoints()...)
	}
}
