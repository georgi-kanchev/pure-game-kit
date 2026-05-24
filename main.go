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

	textbox.Color = palette.DarkGray

	textbox.Effects = graphics.NewEffects()
	textbox.Effects.TextLineHeight = 100
	textbox.Effects.TextShadowBlur = 64
	textbox.Effects.TextShadowOffsetX = 80
	textbox.Effects.TextShadowOffsetY = 50
	textbox.Effects.OutlineSize = 64
	textbox.Effects.OutlineColor = palette.Red
	textbox.Roundness = 0.5
	textbox.Effects.BorderSize = 20
	textbox.Effects.BorderColor = palette.Red

	window.SetTargetFPS(60)

	var img = assets.LoadImage("examples/data/desert-0.png")
	var sprite = graphics.NewImage(0, 0, 3, img)
	img.SetCrop(0, 0, 200, 200)
	sprite.Roundness = 0.5
	sprite.Effects = graphics.NewEffects()
	sprite.Effects.BorderSize = 20
	sprite.Effects.BorderColor = palette.Red

	var shape = graphics.NewShapeRoundedRectangle(0, 0, 1000, 500, 0, 0.5, palette.Red)
	_ = shape

	for window.KeepOpen() {
		textbox.Text = debug.MemoryUsage()

		// obj.Effects.TextShadowOffsetX = int8(number.Map(number.Sine(time.Running()/2), -1, 1, -128, 127))
		view.DrawObjects(&sprite)
	}
}
