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
	var obj = graphics.NewTextbox(0, 0, 1500, 1500, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")

	obj.Color = palette.DarkGray

	obj.Effects = graphics.NewEffects()
	obj.Effects.TextLineHeight = 100
	obj.Effects.TextShadowBlur = 64
	obj.Effects.TextShadowOffsetX = 80
	obj.Effects.TextShadowOffsetY = 50
	obj.Effects.OutlineSize = 64
	obj.Effects.OutlineColor = palette.Red

	window.SetTargetFPS(60)

	var obj2 = graphics.NewShapeCircle(0, 0, 500, palette.Red)
	obj2.Roundness = 0.5

	for window.KeepOpen() {
		obj.Text = debug.MemoryUsage()

		// obj.Effects.TextShadowOffsetX = int8(number.Map(number.Sine(time.Running()/2), -1, 1, -128, 127))
		view.DrawObjects(&obj, &obj2)
	}
}
