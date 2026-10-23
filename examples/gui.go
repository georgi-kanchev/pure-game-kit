package example

import (
	"fmt"
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
	var layout = assets.LoadGUILayout("tools/ui-layout-editor/test-layout.xml")
	var theme = assets.LoadGUITheme("examples/data/theme.xml")
	var view = graphics.NewView(1)

	_ = theme

	var boxCols, itemCols = [5]uint{}, [20]uint{}
	for i := range boxCols {
		boxCols[i] = color.RandomDark()
	}
	for i := range itemCols {
		itemCols[i] = color.RandomDark()
	}

	window.SetTargetFPS(0)

	gui.Scale = 1

	var menuItems = []string{"item1", "item2", "item3", "item4", "item5"}
	var hor, ver float32
	var hor2, ver2 float32
	var input = "hello, world! test"
	var input2 = "tttttt"
	var s float32 = 0.5
	for window.KeepOpen() {
		for i, c := range boxCols {
			var area, _, _ = layout.Box(i, 1)
			gui.Object(0, 0, 0, 0, c, area, gui.Area{}, false)
		}

		for i, c := range itemCols {
			var area, mask = layout.Item(i, hor, ver, 1)
			gui.Object(0, 0, 0, 0, c, area, mask, false)

			switch i {
			case 0:
				gui.Text("Victory", area, mask, 0, false, false)
			case 1:
				gui.Text("(4 rounds)", area, mask, 0, false, false)
			case 5:
				gui.Text("UNIT\ntest\nhi", area, mask, 0, false, false)
			case 6:
				gui.Inputbox(&input2, "mhm...", area, mask, 0, true)
			}
		}

		var area = gui.AreaHUD(0.5, 1, 700, 100)
		area.Y -= 50

		gui.Inputbox(&input, "enter name...", area, gui.Area{}, 0, true)
		gui.Button("button", geometry.NewArea(0, 0, 200, 50), gui.Area{}, 0, true)
		gui.Slider(&s, 0, geometry.NewArea(0, 200, 200, 50), gui.Area{}, 0, true)
		gui.Image(0, geometry.NewArea(-200, -200, 100, 100), true)

		var unitsArea, ucw, uch = layout.Box(3, 1)
		gui.Scrolls(&hor, &ver, ucw, uch, unitsArea, 0)

		if !gui.IsAnyTyping() && keyboard.IsKeyPressed(key.A) {
			var area, cw, ch = layout.Box(5, 1)
			gui.Object(0, 0, 0, 0, palette.Azure, area, gui.Area{}, true)
			gui.Scrolls(&hor2, &ver2, cw, ch, area, 0)
			for i := range 4 {
				var area, mask = layout.Item(20+i, hor2, ver2, 1)
				gui.Object(0, 0, 0, 0, palette.Beige, area, mask, true)
			}
		}

		gui.Menu(geometry.NewArea(1100, 0, 200, 200), gui.Area{}, 0, true, 60, menuItems...)
		for i := range 5 {
			if gui.IsMenuJustClicked(i) {
				fmt.Printf("menuItems[i]: %v\n", menuItems[i])
			}
		}

		gui.Menu(geometry.NewArea(-1100, 0, 200, 200), gui.Area{}, 0, true, 60, menuItems...)
		for i := range 5 {
			if gui.IsMenuJustClicked(i) {
				fmt.Printf("%v\n", menuItems[i])
			}
		}

		view.DrawDebugInfo(false)
	}
}
