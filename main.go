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
	textbox.Effects = graphics.NewEffects()
	textbox.Effects.FillColor = palette.DarkGray

	window.SetTargetFPS(60)

	var img = assets.LoadImage("examples/data/desert-0.png")
	var sprite = graphics.NewImage(0, 0, 3, img)
	sprite.Roundness = 0.5

	var shape = graphics.NewShapeRoundedRectangle(0, 0, 1000, 500, 0, 0.5)
	_ = shape

	for window.KeepOpen() {
		textbox.Text = debug.MemoryUsage()

		// sprite.Width = number.Map(number.Sine(time.Running()/2), -1, 1, 500, 2500)
		// sprite.Effects.BorderSize = number.Map(number.Sine(time.Running()/2), -1, 1, -300, 300)
		// sprite.Roundness = number.Map(number.Sine(time.Running()), -1, 1, 0, 1)
		// sprite.ImageCropArea.X, sprite.ImageCropArea.Y = view.MousePosition()
		// sprite.ImageCropArea.Width, sprite.ImageCropArea.Height = img.Size()
		view.DrawObjects(&textbox)
	}
}
