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
	var animations = assets.LoadAnimations(units, "examples/data/animations.xml")
	var idle = motion.NewAnimationFromAsset(animations, "man-idle", 6, true)
	var walk = motion.NewAnimationFromAsset(animations, "man-walk", 8, true)

	for window.KeepOpen() {
		var frame = idle.Frame()
		if keyboard.IsAnyKeyPressed() {
			frame = walk.Frame()
		}
		view.DrawImage(0, 0, 16*2, 16*3, 0, frame, palette.White, geometry.Area{})
	}
}
