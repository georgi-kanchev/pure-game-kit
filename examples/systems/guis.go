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
		gui.Container("top", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "200",
			[][2]string{
				{property.RGBA, "255 0 0 128"},
			}),
		gui.NewButton(btn1, d.OwnerLeftX+"+20", d.OwnerTopY+"+20", "200", d.OwnerHeight+"",
			[][2]string{
				{property.RGBA, "0 255 0 255"},
			},
			[][2]string{
				{property.Text, "Hello,\nWorld!\nNew\nLine"},
			}),
	)

	cam.Angle = 45

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.ScreenWidth /= 2
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		menu.Draw(&cam)
	}
}
