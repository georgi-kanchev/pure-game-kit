package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
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
	textbox.Effects.TextLineHeight = 100
	textbox.Effects.TextBackColor = palette.Red

	window.SetTargetFPS(60)

	for window.KeepOpen() {
		textbox.Text = debug.MemoryUsage()

		view.DrawObjects(&textbox)
	}
}
