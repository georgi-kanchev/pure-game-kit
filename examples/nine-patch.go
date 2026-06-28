package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func NinePatch() {
	window.Create("example - nine patch", false, true)
	var view = graphics.NewView(10)
	var img = assets.LoadImage("examples/data/9patch2.png")
	var ninePatch = assets.LoadImageCrop9Patch(img, 15, 15, 15, 15)

	ninePatch = img

	for window.KeepOpen() {
		var sine = number.Sine(time.Running())
		view.DrawImage(0, 0, sine*100, sine*50, 0, ninePatch, palette.White, geometry.Area{})
	}
}
