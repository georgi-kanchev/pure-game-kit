package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/window"
)

func main() {
	var view = graphics.NewView(1)
	var font = assets.LoadFont2("tools/sdf-font-generator/results/Montserrat-Medium.png",
		"tools/sdf-font-generator/results/Montserrat-Medium.xml")
	var tb = graphics.NewTextbox(font, 0, 0)
	tb.Text = "Hello, World!"

	assets.LoadDefaultFont()

	for window.KeepOpen() {
		window.Title = "pure-game-kit: hub"
		window.TargetFPS = 0

		view.DrawTextboxes(tb)
		view.DrawTextDebug(true, true, false, true)
	}
}
