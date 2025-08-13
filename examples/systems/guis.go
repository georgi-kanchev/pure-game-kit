package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GUIs() {
	var cam = graphics.NewCamera(1)
	var menu = gui.New(
		gui.Container("top", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "600",
			p.RGBA, "255 0 0 128"),
		gui.NewButton("btn1", "200", "100", p.RGBA, "0 255 0 255", p.OffsetX, "10", p.OffsetY, "10"),
		gui.NewButton("btn2", "100", "150", p.RGBA, "255 255 0 255"),
		gui.NewButton("btn3", "150", "180", p.RGBA, "0 255 255 255"),
		gui.NewButton("btn4", "100", "180", p.RGBA, "128 128 255 255", p.NewRow, ""),
		gui.NewButton("btn5", "200", "80", p.RGBA, "255 0 255 255", p.OffsetX, "20"),
		gui.NewButton("btn6", "150", "150", p.RGBA, "255 128 128 255"),
		gui.NewButton("btn7", "800", "350", p.RGBA, "0 0 0 255",
			p.Text, "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
	)

	cam.Angle = 45

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		menu.Draw(&cam)
	}
}
