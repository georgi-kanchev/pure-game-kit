package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", true, true)
	var view = graphics.NewView(1)

	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")
	_ = layout

	for window.KeepOpen() {
		view.DrawColor(palette.Black)
	}
}
