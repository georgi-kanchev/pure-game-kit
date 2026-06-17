package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", false, true)
	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")

	var boxCols, itemCols = [5]uint{}, [20]uint{}
	for i := range boxCols {
		boxCols[i] = color.RandomDark()
	}
	for i := range itemCols {
		itemCols[i] = color.RandomDark()
	}

	window.SetTargetFPS(0)

	layout.SetVisibleItem(2, false)

	var a float32
	for window.KeepOpen() {
		for i, c := range boxCols {
			var area, _, _ = layout.Box(i)
			gui.Shape(c, 0, area, assets.Area{})
		}

		a = number.Map(number.Sine(time.Running()), -1, 1, 0, 1)
		gui.Scale = 0.5 + a/2

		for i, c := range itemCols {
			var area, mask = layout.Item(i, 0, a)
			gui.Shape(c, 0, area, mask)
			switch i {
			case 0:
				gui.Label("Victory", area, mask)
			case 1:
				gui.Label("(4 rounds)", area, mask)
			case 5:
				gui.Label("UNIT", area, mask)
			}
		}
		var area = gui.AreaHUD(0.5, 1, 0.2, 200)
		area.Y -= 50
		gui.Shape(palette.Red, 1, area, assets.Area{})
	}
}
