package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GUIs() {
	var cam = graphics.NewCamera(1)
	// var font = assets.LoadFonts(32, "examples/data/monogram.ttf")[0]
	var menu = gui.New(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("button", p.TextFontId, "", p.TextLineGap, "-0.2", p.TextLineHeight, "30",
			p.Color, "128 128 255 255", p.Width, "300", p.Height, "100"),
		gui.Container("top", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "600",
			p.Color, "255 0 0 128", p.ThemeId, "button"),
		// gui.WidgetButton("btn1", p.OffsetX, "10", p.OffsetY, "10"),
		// gui.WidgetButton("btn2"),
		// gui.WidgetButton("btn3"),
		// gui.WidgetButton("btn4", p.NewRow, ""),
		// gui.WidgetButton("btn5"),
		// gui.WidgetButton("btn6"),
		gui.WidgetButton("btn7", p.Width, "2000", p.Height, "300",
			p.Text, "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
	)

	cam.Angle = 45

	assets.LoadDefaultFont()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		menu.Draw(&cam)
	}
}
