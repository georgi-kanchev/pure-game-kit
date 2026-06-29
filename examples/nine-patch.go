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
	var img = assets.LoadImage("examples/data/Preview.png")
	var crop = assets.LoadImageCrop(img, 523, 208, 116, 39)
	var ninePatch = assets.LoadImage9Patch(crop, 10, 10, 10, 10)

	// ninePatch = img

	for window.KeepOpen() {
		var sine = number.Sine(time.Running() / 2)
		view.DrawImage(0, 0, sine*100, sine*100, 45, ninePatch, palette.White, geometry.Area{})
	}
}
