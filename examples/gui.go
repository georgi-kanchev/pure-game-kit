package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
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

	gui.Scale = 1

	var hor, ver float32
	var hor2, ver2 float32
	var value float32
	for window.KeepOpen() {
		for i, c := range boxCols {
			var area, _, _ = layout.Box(i)
			gui.Shape(c, 0, area, geometry.Area{})
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

		var area = gui.AreaHUD(0.5, 1, 700, 100)
		area.X += 50
		area.Y -= 50
		gui.Slider(&value, 0.1, area, geometry.Area{})

		if keyboard.IsKeyPressed(key.A) {
			var area, cw, ch = layout.Box(5)
			gui.Shape(palette.Azure, 0, area, geometry.Area{})
			gui.Scrolls(&hor2, &ver2, cw, ch, area)
			for i := range 4 {
				var area, mask = layout.Item(20+i, hor2, ver2)
				gui.Shape(palette.Beige, 0, area, mask)
			}
		}

		var unitsArea, ucw, uch = layout.Box(3)
		gui.Scrolls(&hor, &ver, ucw, uch, unitsArea)

		view.DrawDebugInfo(false)
	}
}
