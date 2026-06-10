package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func NinePatch() {
	window.Create("example - nine-patch", false, true)

	var img = assets.LoadImage("examples/data/ribbons.png")
	var crop = assets.LoadImageCrop(img, 0, 0, 48, 32)
	var ninePatch = assets.LoadNinePatch(crop, 16, 16, 16, 16, false, false, false, false, false)
	var view = graphics.NewView(1)
	var obj = graphics.NewNinePatch(0, 0, 300, 200, ninePatch)

	for window.KeepOpen() {
		var s = number.Map(number.Sine(time.Running()/2), -1, 1, 0.5, 1.5)
		obj.Width, obj.Height = 300*s, 200*s

		view.DrawObject(&obj)
	}
}
