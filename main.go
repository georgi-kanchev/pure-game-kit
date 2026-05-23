package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	// var font = assets.LoadFont("tools/msdf-atlas-gen/font.png", "tools/msdf-atlas-gen/font.json")
	var obj = graphics.NewTextbox(0, 0, 400, 400, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")

	obj.Width *= 4
	obj.Height *= 4

	obj.Effects = graphics.NewEffects()

	obj.Color = palette.DarkGray

	// obj.Mask = graphics.NewArea(obj.X-obj.Width/2, obj.Y-obj.Height/2, obj.Width, obj.Height)
	obj.Roundness = 0.5
	obj.TextOutlineSize = 0
	obj.TextOutlineColor = palette.Black
	obj.TextWeight = 128

	// window.SetTargetFPS(60)

	obj.TextShadowBlur = 50
	obj.TextShadowOffsetX = 0
	obj.TextShadowOffsetY = -127
	obj.TextShadowSize = 100
	obj.TextShadowColor = palette.Black

	for window.KeepOpen() {
		obj.Text = debug.MemoryUsage()
		obj.TextLineHeight = number.Map(number.Sine(time.Running()/2), -1, 1, 0, 255)
		// obj.Mask.X, obj.Mask.Y = view.MousePosition()
		// obj.Mask.X -= obj.Mask.Width / 2
		// obj.Mask.Y -= obj.Mask.Height / 2

		view.DrawObjects(&obj)
	}
}
