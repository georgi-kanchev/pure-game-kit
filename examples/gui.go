package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func GUI() {
	window.Create("example - gui", false, true)
	var layout = assets.LoadLayout("tools/ui-layout-editor/test-layout.xml")
	var view = graphics.NewView(1)

	var boxCols, itemCols = [5]uint{}, [20]uint{}
	for i := range boxCols {
		boxCols[i] = color.RandomDark()
	}
	for i := range itemCols {
		itemCols[i] = color.RandomDark()
	}

	// window.SetTargetFPS(0)

	var hor, ver float32
	var hor2, ver2 float32
	var input = "hello, world!"
	for window.KeepOpen() {
		for i, c := range boxCols {
			gui.Shape(c, 0, layout.Box(i), geometry.Area{})
		}

		for i, c := range itemCols {
			var area, mask = layout.Item(i, hor, ver)
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
		var area = gui.AreaHUD(0.5, 1, 450, 150)
		area.X += 50
		area.Y -= 50
		gui.Inputbox(&input, area, geometry.Area{})

		gui.Scrolls(layout, 3, &hor, &ver)

		gui.Shape(palette.Azure, 0, layout.Box(5), geometry.Area{})
		gui.Scrolls(layout, 5, &hor2, &ver2)
		for i := range 4 {
			var area, mask = layout.Item(20+i, hor2, ver2)
			gui.Shape(palette.Beige, 0, area, mask)
		}

		view.DrawDebugInfo(false)
	}
}
