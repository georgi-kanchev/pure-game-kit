package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", true, true)
	var view = graphics.NewView(1)

	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")

	var boxCols, itemCols = [4]uint{}, [19]uint{}
	for i := range boxCols {
		boxCols[i] = color.RandomDark()
	}
	for i := range itemCols {
		itemCols[i] = color.RandomDark()
	}

	for window.KeepOpen() {
		for i, c := range boxCols {
			var x, y, w, h = layout.Box(i, view.Zoom)
			view.DrawShape(x, y, w, h, 0, 0, c)
		}

		// for i, c := range itemCols {
		// 	var x, y, w, h = layout.Item(i, view.Zoom)
		// 	view.DrawShape(x, y, w, h, 0, 0, c)
		// }
	}
}
