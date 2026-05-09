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
	var obj = graphics.NewObject(0, 0)
	obj.TextFont = font
	obj.Text = "Hello, World!"

	assets.LoadDefaultFont()

	for window.KeepOpen() {
		window.Title = "pure-game-kit: hub"
		window.TargetFPS = 0

		view.DrawObjects(&obj)
		view.DrawTextDebug(true, true, false, true)
	}
}
