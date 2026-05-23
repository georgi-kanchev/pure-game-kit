package main

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	// var font = assets.LoadFont("tools/msdf-atlas-gen/font.png", "tools/msdf-atlas-gen/font.json")
	var obj = graphics.NewTextbox(0, 0, 300, 300, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")

	obj.Width *= 4
	obj.Height *= 4
	obj.TextLineHeight = 40

	obj.Effects = graphics.NewEffects()
	obj.Effects.OutlineSize = 255
	obj.Effects.OutlineColor = palette.Red

	obj.Color = palette.DarkGray

	obj.Mask = graphics.NewArea(obj.X-obj.Width/2, obj.Y-obj.Height/2, obj.Width, obj.Height)
	// obj.Angle = 1
	obj.Roundness = 0.5

	var obj2 = graphics.NewImage(0, 0, assets.LoadImage("examples/data/flail.PNG"))
	obj2.Width *= 3
	obj2.Height *= 3

	// window.SetTargetFPS(60)

	obj.Angle = 20

	fmt.Printf("debug.LinesOfCode(): %v\n", debug.LinesOfCode())

	var drawObj = false
	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			drawObj = true
		}
		obj.Text = debug.MemoryUsage()
		// obj.TextWeight = byte(number.Map(number.Sine(time.Running()), -1, 1, 0, 255))
		obj.Mask.X, obj.Mask.Y = view.MousePosition()
		obj.Mask.X -= obj.Mask.Width / 2
		obj.Mask.Y -= obj.Mask.Height / 2
		obj2.X, obj2.Y = view.MousePosition()
		obj2.Angle += time.Delta() * 5

		view.DrawObjects(&obj)
		if drawObj {
			view.DrawObjects(&obj2)
		}
	}
}
