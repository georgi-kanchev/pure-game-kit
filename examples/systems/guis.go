package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

const btn1 = "btn1"

func GUIs() {
	var cam = graphics.NewCamera(1)
	var menu = gui.New(
		gui.Container("top", d.CameraLeftX+"+50", d.CameraTopY+"+50", d.CameraWidth+"-100", "200",
			property.RGBA, "255 0 0 128"),
		gui.NewButton(btn1, d.OwnerLeftX+"+10", d.OwnerTopY+"+10", "300", d.OwnerHeight+"-20",
			property.RGBA, "0 255 0 255"),
	)

	cam.Angle = 45

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		menu.Draw(&cam)
	}
}
