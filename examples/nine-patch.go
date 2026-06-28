package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func NinePatch() {
	window.Create("example - nine patch", false, true)
	var view = graphics.NewView(1)

	for window.KeepOpen() {
		view.DrawColor(palette.Red)
	}
}
