package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Events() {
	window.Create("game", true, true)
	var view = graphics.NewView(1)
	var img = assets.LoadImage("examples/data/flail.PNG")
	// img.Unload()
	for window.KeepOpen() {
		view.DrawImage(0, 0, 200, 200, 0, img, palette.White)
	}
}
