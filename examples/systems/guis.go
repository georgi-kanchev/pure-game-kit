package example

import (
	"fmt"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/window"
)

const btn1 = "btn1"

func GUIs() {
	var cam = graphics.NewCamera(1)
	var menu = gui.New(
		gui.NewButton(btn1, d.OwnerLeft, d.OwnerRight, d.MyTextWidth, d.OwnerHeight),
	)

	menu.SetProperty(btn1, p.X, "test")

	fmt.Printf("%v\n", menu.Property(btn1, p.X))

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		menu.Draw(&cam)
	}
}
