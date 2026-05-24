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
	var obj = graphics.NewTextbox(0, 0, 2000, 1500, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")

	obj.Color = palette.DarkGray

	obj.Effects = graphics.NewEffects()
	obj.Effects.TextLineHeight = 100
	obj.Effects.TextShadowBlur = 64
	obj.Effects.TextShadowOffsetX = 80
	obj.Effects.TextShadowOffsetY = 50
	obj.Effects.OutlineSize = 64
	obj.Effects.OutlineColor = palette.Red
	obj.Roundness = 0.5
	obj.Effects.BorderSize = 20
	obj.Effects.BorderColor = palette.Red

	window.SetTargetFPS(60)

	var img = assets.LoadImage("examples/data/desert-0.png")
	var obj2 = graphics.NewImage(0, 0, 4, img)
	// obj2.ImageCropArea = graphics.NewArea(0, 0, 300, 300)
	obj2.Roundness = 0.5

	obj2.Effects = graphics.NewEffects()
	obj2.Effects.BorderSize = 10
	obj2.Effects.BorderColor = palette.Red
	// obj2.Width = 500

	for window.KeepOpen() {
		obj.Text = debug.MemoryUsage()

		// obj.Effects.TextShadowOffsetX = int8(number.Map(number.Sine(time.Running()/2), -1, 1, -128, 127))
		view.DrawObjects(&obj2)
	}
}
