package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/gui"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
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

	gui.Scale = 1.2

	var hor, ver float32
	var hor2, ver2 float32
	var input = "hello, world! test"
	var input2 = "tttttt"
	var s float32 = 0.5
	for window.KeepOpen() {
		for i, c := range boxCols {
			var area, _, _ = layout.Box(i)
			gui.Shape(0, c, area, gui.Area{}, 0, false)
		}

		if mouse.ScrollY() != 0 {
			print()
		}

		for i, c := range itemCols {
			var area, mask = layout.Item(i, hor, ver)
			gui.Shape(0, c, area, mask, 0, false)

			switch i {
			case 0:
				gui.Label("Victory", area, mask, 0, false)
			case 1:
				gui.Label("(4 rounds)", area, mask, 0, false)
			case 5:
				gui.Text("UNIT\ntest\nhi", 50, area, mask, 0, false)
			case 6:
				gui.Inputbox(&input2, "mhm...", area, mask, 0, true)
			}
		}

		var unitsArea, ucw, uch = layout.Box(3)
		gui.Scrolls(&hor, &ver, ucw, uch, unitsArea, 0)

		if !gui.IsAnyTyping() && keyboard.IsKeyPressed(key.A) {
			var area, cw, ch = layout.Box(5)
			gui.Shape(0, palette.Azure, area, gui.Area{}, 0, false)
			gui.Scrolls(&hor2, &ver2, cw, ch, area, 0)
			for i := range 4 {
				var area, mask = layout.Item(20+i, hor2, ver2)
				gui.Shape(0, palette.Beige, area, mask, 0, false)
			}
		}

		var area = gui.AreaHUD(0.5, 1, 700, 100)
		area.Y -= 50

		gui.Inputbox(&input, "enter name...", area, gui.Area{}, 0, true)

		gui.Button(geometry.NewArea(0, 0, 200, 50), gui.Area{}, 0, true)
		gui.Label("button", geometry.NewArea(0, 0, 200, 50), gui.Area{}, 0, false)

		gui.Slider(&s, 0, geometry.NewArea(0, 200, 200, 50), gui.Area{}, 0, true)

		view.DrawDebugInfo(true)
	}
}
