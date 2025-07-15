package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Texts() {
	var cam = graphics.NewCamera(3)
	var f = assets.LoadFonts(24, "font.ttf")[0]
	var node = graphics.NewNode("font")

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawText(f, "кирилица", -400, -300, 1000, 0, 0.4, 0.7, color.Red)
		cam.DrawText(f, "кирилица", -400, -300, 1000, 0, 0, 0.5, color.White)

		cam.DrawNodes(&node)
	}
}
