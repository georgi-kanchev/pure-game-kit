package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/graphics"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 0, false, false)

	var view = graphics.NewView(1)
	var path = "tools/sdf-font-generator/results/Roboto-Bold."
	var font = assets.LoadFont2(path+"png", path+"xml")
	var obj = graphics.NewObject(0, 0)
	obj.TextFont = font
	obj.Text = "Hello, World!"

	// assets.LoadDefaultFont()
	engine.Run(func() {
		view.DrawObjects(&obj)
	})
}
