package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Texts() {
	var cam = graphics.NewCamera(3)
	var f = assets.LoadFonts(32, "font.ttf")[0]
	var textBox = graphics.NewTextBox(f, 0, 0, "test: ", 145, " hello: ", 851.32, "...")
	// textBox.FontId = ""
	// textBox.ValueScale = 5

	var shadow = textBox
	shadow.Thickness = 0.8
	shadow.Color = color.Black

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(1, 32, color.Darken(color.Gray, 0.5))

		var x, y, w, h = textBox.Rectangle()
		cam.DrawRectangle(x, y, w, h, textBox.Angle, color.Red)
		cam.DrawTextBoxes(&shadow, &textBox)
	}
}
