package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", false, true)
	var view = graphics.NewView(1)
	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")

	var boxCols, itemCols = [5]uint{}, [20]uint{}
	for i := range boxCols {
		boxCols[i] = color.RandomDark()
	}
	for i := range itemCols {
		itemCols[i] = color.RandomDark()
	}

	window.SetTargetFPS(0)

	for window.KeepOpen() {
		for i, c := range boxCols {
			var x, y, w, h = layout.BoxArea(i, view.Zoom)
			view.DrawShape(x, y, w, h, 0, 0, c, graphics.Area{})
		}

		for i, c := range itemCols {
			var x, y, w, h = layout.ItemArea(i, view.Zoom, 0, 1)
			var mask = graphics.NewArea(layout.ItemMask(i, view.Zoom))
			view.DrawShape(x, y, w, h, 0, 0, c, mask)
		}
	}
}
