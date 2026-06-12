package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", true, true)
	var view = graphics.NewView(1.2)

	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")

	for window.KeepOpen() {
		var x, y, w, h = layout.Box(0, view.Zoom)
		view.DrawShape(x, y, w, h, 0, 0, palette.White)

		var x1, y1, w1, h1 = layout.Box(1, view.Zoom)
		view.DrawShape(x1, y1, w1, h1, 0, 0, palette.Red)
	}
}
