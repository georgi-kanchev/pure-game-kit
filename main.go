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
	var textbox = graphics.NewTextbox(-200, 0, 2000, 1500, font, "((")
	textbox.Effects.FillColor = palette.DarkGray
	textbox.Effects.TextUnderline = true
	textbox.Angle = 5

	window.SetTargetFPS(60)

	for window.KeepOpen() {
		textbox.Text = debug.MemoryUsage()
		var x, _ = view.MousePosition()
		textbox.Effects.TextLineHeight = 70 + x/10

		view.DrawObjects(&textbox)
	}
}
