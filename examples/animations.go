package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Animations() {
	window.Create("examples - animation", true, true)
	var view = graphics.NewView(10)
	var units = assets.LoadImage("examples/data/units.png")
	var atlas = assets.LoadAtlas(units, "examples/data/animations.xml")
	var idle = motion.NewAnimation(6, true, atlas.Crops("man-idle")...)
	var walk = motion.NewAnimation(8, true, atlas.Crops("man-walk")...)

	for window.KeepOpen() {
		var frame = idle.Frame()
		if keyboard.IsAnyKeyPressed() {
			frame = walk.Frame()
		}
		view.DrawImage(0, 0, 16*2, 16*3, 0, frame, palette.White, geometry.Area{})
	}
}
