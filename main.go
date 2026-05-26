package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	// var font = assets.LoadFont("tools/msdf-atlas-gen/font.png", "tools/msdf-atlas-gen/font.json")
	var textbox = graphics.NewTextbox(0, 0, 2000, 1500, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")
	textbox.Effects.FillColor = palette.DarkGray
	textbox.Effects.TextLineHeight = 40
	textbox.Effects.TextBackColor = palette.Red

	// window.SetTargetFPS(60)

	var rect = graphics.NewShapeRectangle(0, 0, 0, 0, 0)
	rect.Effects.BorderSize = 0
	rect.Effects.FillColor = color.RGBA(255, 255, 255, 100)

	for window.KeepOpen() {
		textbox.Text = debug.MemoryUsage()

		var mx, my = view.MousePosition()
		var w, h = textbox.TextMeasureLine(0, 99999)
		rect.X, rect.Y = mx, my
		rect.Width, rect.Height = w, h

		view.DrawObjects(&textbox, &rect)
	}
}
