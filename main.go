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
	var obj = graphics.NewTextbox(0, 0, 1000, 1000, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")

	obj.Effects = graphics.NewEffects()
	obj.Color = palette.DarkGray

	obj.Effects.TextShadowSize = 128
	obj.Effects.TextShadowOffsetX = 10

	window.SetTargetFPS(60)

	for window.KeepOpen() {
		obj.Text = debug.MemoryUsage()
		// obj.TextShadowOffsetX = int8(number.Map(number.Sine(time.Running()/2), -1, 1, 0, 127))
		// obj.Mask.X, obj.Mask.Y = view.MousePosition()
		// obj.Mask.X -= obj.Mask.Width / 2
		// obj.Mask.Y -= obj.Mask.Height / 2

		view.DrawObjects(&obj)
	}
}
