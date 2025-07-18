package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var f = assets.LoadFonts(32, "font.ttf")[0]
	var textBox = graphics.NewTextBox(f, 0, 0, "test: ", 145, " hello: ", 851.32)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.Color = color.Red
		cam.DrawNodes(&textBox.Node)
		textBox.Color = color.White
		cam.DrawTextBoxes(&textBox)
	}
}
